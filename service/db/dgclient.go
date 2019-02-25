package db

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/hashicorp/golang-lru"
	"github.com/intuit/katlas/service/metrics"
	"github.com/intuit/katlas/service/util"
	"google.golang.org/grpc"
)

// Action as oper
type Action int

const (
	create Action = iota
	update
	delete
)

//CacheKey - Define key name for LruCache
const CacheKey = "dbSchema"

//LruCache - Define type LRU Cache
var LruCache *lru.Cache

//InitLruCacheDBSchema - a flag to indicate if LruCache has initial DBSchema when server starts
var InitLruCacheDBSchema bool

// Schema dgraph database schema
type Schema struct {
	Predicate string   `json:"predicate"`
	Type      string   `json:"type"`
	List      bool     `json:"list,omitempty"`
	Index     bool     `json:"index,omitempty"`
	Upsert    bool     `json:"upsert,omitempty"`
	Count     bool     `json:"count,omitempty"`
	Reverse   bool     `json:"reverse,omitempty"`
	Tokenizer []string `json:"tokenizer,omitempty"`
}

// DGClient will run query or command on dgraph
type DGClient struct {
	conn *grpc.ClientConn
	dc   *dgo.Dgraph
}

// IDGClient ... define interface to DGClient
type IDGClient interface {
	GetCacheContainsDBSchema() (*lru.Cache, error)
	GetSchemaFromCache(cache *lru.Cache) ([]*api.SchemaNode, error)
	RemoveDBSchemaFromCache(cache *lru.Cache)
	GetSchemaFromDB() ([]*api.SchemaNode, error)
	CreateSchema(sm Schema) error
	DropSchema(name string) error
	GetEntity(uuid string) (map[string]interface{}, error)
	GetAllByClusterAndType(meta string, cluster string) (map[string]interface{}, error)
	DeleteEntity(uuid string) error
	CreateEntity(meta string, data map[string]interface{}) (string, error)
	CreateOrDeleteEdge(fromType string, fromUID string, toType string, toUID string, rel string, op Action) error
	SetFieldToNull(delMap map[string]interface{}) error
	UpdateEntity(uuid string, data map[string]interface{}) error
	GetQueryResult(query string) (map[string]interface{}, error)
	Close() error
	ExecuteDgraphQuery(query string) (map[string]interface{}, error)
}

// NewDGClient create client instance
// TODO:
// consider to return single client stub without close connection
func NewDGClient(dgraphHost string) *DGClient {
	// Dial a gRPC connection.
	log.Infof("Connecting to dgraph [%s]", dgraphHost)
	conn, err := grpc.Dial(dgraphHost,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(20*1024*1024)),
		grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}
	return &DGClient{
		conn, dgo.NewDgraphClient(api.NewDgraphClient(conn)),
	}
}

// GetEntity - get entity by uid
func (s DGClient) GetEntity(uuid string) (map[string]interface{}, error) {
	q := `
		query qry($uuid: string) {
			objects(func: uid($uuid)) {
				uid
				expand(_all_) {
					uid
					expand(_all_)
				}
			}
		}
	`
	resp, err := s.dc.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uuid": uuid})
	if err != nil {
		metrics.DgraphNumGetEntityErr.Inc()
		metrics.DgraphNumQueriesErr.Inc()
		return nil, err
	}
	metrics.DgraphNumQueries.Inc()

	m := make(map[string]interface{})
	err = json.Unmarshal(resp.Json, &m)
	if err != nil {
		return nil, err
	}
	if len(m[util.Objects].([]interface{})) > 0 {
		// only uid return, means no record found
		data := m[util.Objects].([]interface{})[0].(map[string]interface{})
		if _, ok := data[util.UID]; ok && len(data) == 1 {
			return map[string]interface{}{}, nil
		}
	}
	return m, nil
}

// DeleteEntity - delete entity by uuid
func (s DGClient) DeleteEntity(uuid string) error {
	ctx := context.Background()
	txn := s.dc.NewTxn()
	defer txn.Discard(ctx)
	q := `
		{
  			"uid": "` + uuid + `"
		}
    `
	mu := &api.Mutation{
		CommitNow:  true,
		DeleteJson: []byte(q),
	}
	_, err := txn.Mutate(ctx, mu)
	if err != nil {
		metrics.DgraphNumDeleteEntityErr.Inc()
		metrics.DgraphNumMutationsErr.Inc()
		log.Debug(err)
		return err
	}
	metrics.DgraphNumMutations.Inc()
	return nil
}

// CreateEntity - create entity
func (s DGClient) CreateEntity(meta string, data map[string]interface{}) (string, error) {
	ctx := context.Background()
	txn := s.dc.NewTxn()
	defer txn.Discard(ctx)
	mu := &api.Mutation{
		CommitNow: true,
	}
	jsonData, _ := json.Marshal(data)
	mu.SetJson = jsonData
	resp, err := txn.Mutate(ctx, mu)
	if err != nil {
		metrics.DgraphNumCreateEntityErr.Inc()
		metrics.DgraphNumMutationsErr.Inc()
		log.Error(err, data)
		return "", err
	}
	metrics.DgraphNumMutations.Inc()

	log.Infof("%s %s created/updated successfully", meta, data["name"])
	// return created blank node uid
	if uid, ok := resp.Uids["A"]; ok {
		return uid, nil
	}
	if uid, ok := resp.Uids["blank-0"]; ok {
		return uid, nil
	}
	return data[util.UID].(string), nil
}

// SetFieldToNull - remove list or edges from nodes
func (s DGClient) SetFieldToNull(delMap map[string]interface{}) error {
	ctx := context.Background()
	txn := s.dc.NewTxn()
	defer txn.Discard(ctx)
	mu := &api.Mutation{
		CommitNow: true,
	}
	delJSON, _ := json.Marshal(delMap)
	mu.DeleteJson = delJSON
	_, err := txn.Mutate(ctx, mu)
	if err != nil {
		metrics.DgraphNumMutationsErr.Inc()
		log.Info(err)
		return err
	}
	metrics.DgraphNumMutations.Inc()
	return nil
}

// CreateOrDeleteEdge - create or remove edge
func (s DGClient) CreateOrDeleteEdge(fromType string, fromUID string, toType string, toUID string, rel string, op Action) error {
	ctx := context.Background()
	txn := s.dc.NewTxn()
	defer txn.Discard(ctx)
	// construct json string for create/delete edge
	var buffer bytes.Buffer
	buffer.WriteString(`{"uid":"`)
	buffer.WriteString(fromUID)
	buffer.WriteString(`","`)
	buffer.WriteString(rel)
	buffer.WriteString(`": {"uid": "`)
	buffer.WriteString(toUID)
	buffer.WriteString(`"}}`)
	mu := &api.Mutation{
		CommitNow: true,
	}
	switch op {
	case create:
		mu.SetJson = []byte(buffer.String())
	case delete:
		mu.DeleteJson = []byte(buffer.String())
	default:
		log.Debug("No operation found, skip")
		return nil
	}
	_, err := txn.Mutate(ctx, mu)
	if err != nil {
		metrics.DgraphNumMutationsErr.Inc()
		log.Debug(err)
		return err
	}
	metrics.DgraphNumMutations.Inc()
	return nil
}

// UpdateEntity - update entity
func (s DGClient) UpdateEntity(uuid string, data map[string]interface{}) error {
	ctx := context.Background()
	txn := s.dc.NewTxn()
	defer txn.Discard(ctx)
	mu := &api.Mutation{
		CommitNow: true,
	}
	data["uid"] = uuid
	jdata, err := json.Marshal(data)
	if err != nil {
		log.Debug(err)
		return err
	}
	mu.SetJson = jdata
	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		metrics.DgraphNumUpdateEntityErr.Inc()
		metrics.DgraphNumMutationsErr.Inc()
		log.Debug(err)
		return err
	}
	log.Infof("%s updated successfully", data["name"])
	metrics.DgraphNumMutations.Inc()
	return nil
}

// GetQueryResult - get Query Results
func (s DGClient) GetQueryResult(query string) (map[string]interface{}, error) {
	resp, err := s.dc.NewTxn().Query(context.Background(), query)
	if err != nil {
		metrics.DgraphNumQueriesErr.Inc()
		log.Errorf("Query[%v] Error [%v]\n", query, err)
		return nil, err
	}
	metrics.DgraphNumQueries.Inc()

	m := make(map[string]interface{})
	err = json.Unmarshal(resp.Json, &m)
	if err != nil {
		log.Errorf("Query[%v] Error [%v]\n", query, err)
		return nil, err
	}
	return m, nil
}

// GetAllByClusterAndType - query to get result by filter edge
func (s DGClient) GetAllByClusterAndType(meta string, cluster string) (map[string]interface{}, error) {
	q := `
	query qry($type: string, $cluster: string) 
	{
  		objects (func: eq (objtype, $type)) @cascade {
			uid
			name
			resourceid
			cluster @filter (eq(name, $cluster)) {
				name
			}
		}
	}`
	resp, err := s.dc.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$type": meta, "$cluster": cluster})
	if err != nil {
		metrics.DgraphNumQueriesErr.Inc()
		log.Errorf("Query[%v] Error [%v]\n", q, err)
		return nil, err
	}
	metrics.DgraphNumQueries.Inc()

	m := make(map[string]interface{})
	err = json.Unmarshal(resp.Json, &m)
	if err != nil {
		log.Errorf("Query[%v] Error [%v]\n", q, err)
		return nil, err
	}
	return m, nil
}

//GetCacheContainsDBSchema - Get cache which contains db schema
func (s DGClient) GetCacheContainsDBSchema() (*lru.Cache, error) {
	//Add db schema to the cache
	if !InitLruCacheDBSchema {
		dbSchemaNodes, err := s.GetSchemaFromDB()
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, err
		}
		LruCache.Add(CacheKey, dbSchemaNodes)
		InitLruCacheDBSchema = true
	} else {
		//Looks up a key's value from the cache
		_, ok := LruCache.Get(CacheKey)
		if !ok {
			dbSchemaNodes, err := s.GetSchemaFromDB()
			if err != nil {
				log.Errorf("err: %v", err)
				return nil, err
			}
			LruCache.Add(CacheKey, dbSchemaNodes)
		}
	}
	return LruCache, nil
}

//GetSchemaFromCache - Get db schema from cache
func (s DGClient) GetSchemaFromCache(cache *lru.Cache) ([]*api.SchemaNode, error) {
	cache, err := s.GetCacheContainsDBSchema()
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	dbSchemaNodesInterface, ok := cache.Get(CacheKey)
	if !ok {
		log.Errorf("err: %v", err)
		return nil, err
	}

	dbSchemaNodes, ok := dbSchemaNodesInterface.([]*api.SchemaNode)
	return dbSchemaNodes, nil
}

//GetSchemaFromDB - get all predicates
func (s DGClient) GetSchemaFromDB() ([]*api.SchemaNode, error) {
	q := `
		schema {}
	`
	resp, err := s.dc.NewTxn().Query(context.Background(), q)
	if err != nil {
		metrics.DgraphNumQueriesErr.Inc()
		log.Errorf("Query [%v] Error [%v]\n", q, err)
		return nil, err
	}
	metrics.DgraphNumQueries.Inc()
	smn := resp.Schema
	return smn, nil
}

//RemoveDBSchemaFromCache - remove DBSchema key from the Cache
func (s DGClient) RemoveDBSchemaFromCache(cache *lru.Cache) {
	cache.Remove(CacheKey)
}

// CreateSchema - create index
func (s DGClient) CreateSchema(sm Schema) error {
	var buffer bytes.Buffer
	buffer.WriteString(sm.Predicate)
	buffer.WriteString(": ")
	if sm.Type == "password" {
		buffer.WriteString(sm.Type)

	} else if sm.Type == util.UID {
		buffer.WriteString(sm.Type)
		if sm.Count {
			buffer.WriteString(" @count")
		}
		if sm.Reverse {
			buffer.WriteString(" @reverse")
		}
	} else {
		if sm.List {
			buffer.WriteString("[" + sm.Type + "]")
			if sm.Count {
				buffer.WriteString(" @count")
			}
		} else {
			buffer.WriteString(sm.Type)
		}
		if sm.Index {
			buffer.WriteString(" @index(")
			for i, v := range sm.Tokenizer {
				buffer.WriteString(v)
				if i != len(sm.Tokenizer)-1 {
					buffer.WriteString(",")
				}
			}
			buffer.WriteString(")")
		}
		if sm.Upsert {
			buffer.WriteString(" @upsert")
		}
	}
	buffer.WriteString(" .")
	ctx := context.Background()
	err := s.dc.Alter(ctx, &api.Operation{Schema: buffer.String()})
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}

// DropSchema remove db schema by name
func (s DGClient) DropSchema(name string) error {
	ctx := context.Background()
	err := s.dc.Alter(ctx, &api.Operation{DropAttr: name})
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}

// Close - destroy connection
func (s DGClient) Close() error {
	return s.conn.Close()
}

// ExecuteDgraphQuery - Takes a dgraph query as a string and executes on a dgraph instance
func (s DGClient) ExecuteDgraphQuery(query string) (map[string]interface{}, error) {

	txn := s.dc.NewTxn()
	defer txn.Discard(context.Background())

	resp, err := txn.Query(context.Background(), query)
	if err != nil {
		metrics.DgraphNumQueriesErr.Inc()
		log.Errorf("query err: %#v\n", err)
		return nil, errors.New("could not successfully execute query. Please try again later\n" + err.Error())
	}
	metrics.DgraphNumQueries.Inc()

	respjson := map[string]interface{}{}

	err = json.Unmarshal(resp.GetJson(), &respjson)
	if err != nil {
		log.Errorf("unmarshal err: %#v\n", err)
		return nil, errors.New("could not successfully handle data from query. Please try again later")
	}

	// log.Infof("response from executing dgraph query: %#v\n", respjson)
	return respjson, nil

}
