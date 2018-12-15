package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/intuit/katlas/service/util"
	"google.golang.org/grpc"
)

//QueryParamKeyword ...param for keyword query
const QueryParamKeyword = "keyword"

// Action as oper
type Action int

const (
	create Action = iota
	update
	delete
)

// Schema dgraph database schema
type Schema struct {
	Predicate string   `json:"predicate"`
	PType     string   `json:"type"`
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
	CreateSchema(sm Schema) error
	DropSchema(name string) error
	GetEntity(meta string, uuid string) (map[string]interface{}, error)
	GetAllByClusterAndType(meta string, cluster string) (map[string]interface{}, error)
	DeleteEntity(uuid string) error
	CreateEntity(meta string, data map[string]interface{}) (map[string]string, error)
	CreateOrDeleteEdge(fromType string, fromUID string, toType string, toUID string, rel string, op Action) error
	SetFieldToNull(delMap map[string]interface{}) error
	UpdateEntity(meta string, uuid string, data map[string]interface{}) error
	GetQueryResult(queryMap map[string][]string) (map[string]interface{}, error)
	Close() error
}

// NewDGClient create client instance
// TODO:
// consider to return single client stub without close connection
func NewDGClient(dgraphHost string) *DGClient {
	// Dial a gRPC connection.
	log.Infof("Connecting to dgraph [%s]", dgraphHost)
	conn, err := grpc.Dial(dgraphHost, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return &DGClient{
		conn, dgo.NewDgraphClient(api.NewDgraphClient(conn)),
	}
}

// GetEntity - get entity by uid
func (s DGClient) GetEntity(meta string, uuid string) (map[string]interface{}, error) {
	q := `
		{
			objects(func: uid(` + uuid + `)) {
                uid
				expand(_all_) {
                    uid
                    expand(_all_)
                }
			}
		}
	`
	log.Infof("print output query %s\n", q)
	resp, err := s.dc.NewTxn().Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(resp.Json, &m)
	if err != nil {
		return nil, err
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
		log.Debug(err)
		return err
	}
	return nil
}

// CreateEntity - create entity
func (s DGClient) CreateEntity(meta string, data map[string]interface{}) (map[string]string, error) {
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
		log.Error(err)
		return nil, err
	}
	log.Infof("%s %s created/updated successfully", meta, data["name"])
	if uid, ok := data["uid"]; ok {
		return map[string]string{data["name"].(string): uid.(string)}, nil
	}
	return resp.Uids, nil
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
		log.Info(err)
		return err
	}
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
		log.Debug(err)
		return err
	}
	return nil
}

// UpdateEntity - update entity
func (s DGClient) UpdateEntity(meta string, uuid string, data map[string]interface{}) error {
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
		log.Debug(err)
		return err
	}
	return nil
}

// GetQueryResult - get Query Results
func (s DGClient) GetQueryResult(queryMap map[string][]string) (map[string]interface{}, error) {

	var q string

	val, ok := queryMap[QueryParamKeyword]
	if ok {
		if val[0] == "" {
			err := fmt.Errorf("Value not specified for Query Param [%s]", QueryParamKeyword)
			return nil, err
		}
		q = s.GetQueryResultByKeyword(val[0])
	} else {
		if len(queryMap) == 0 {
			err := fmt.Errorf("Query Params not specified")
			return nil, err
		}
		q = s.GetQueryResultByKeyValue(queryMap)
	}

	resp, err := s.dc.NewTxn().Query(context.Background(), q)
	if err != nil {
		log.Errorf("Query[%v] Error [%v]\n", q, err)
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(resp.Json, &m)
	if err != nil {
		log.Errorf("Query[%v] Error [%v]\n", q, err)
		return nil, err
	}
	return m, nil
}

// GetQueryResultByKeyword ...Api to get Query Results by keyword
// Keyword query http://<dgraph ip:port>/v1/query?keyword=webapp
func (s DGClient) GetQueryResultByKeyword(keyword string) (q string) {
	cnt := 0
	q = "{"
	for k, v := range QueryPredicates {
		if v.Type == "string" && v.IsIndex == true && IsExpectedIndexType(v.IndexType, "trigram") {

			filter := "obj" + strconv.Itoa(cnt) + "(func:regexp(" + k + ",/" + keyword + "/i)) {"
			q = q + filter + `
				uid
	            expand(_all_) {
					uid
					expand(_all_)
				}
			}
			`
		}
		cnt++
	}
	q = q + "}"
	return q
}

// GetQueryResultByKeyValue ...Api to get Query Results by key-value pairs
// Key-Value query http://<dgraph ip:port>/v1/query?name=webapp&objtype=Container
func (s DGClient) GetQueryResultByKeyValue(queryMap map[string][]string) (q string) {

	//Only indexed fields can be filtered on
	//Time must be in correct format "2018-10-18 14:36:32 -0700 PDT"
	qps := []string{}
	var funcStr, filterStr string

	for k, v := range queryMap {
		qp := "eq(" + k + ",\"" + v[0] + "\")"
		qps = append(qps, qp)
	}

	funcStr = "(func:" + qps[0] + ") "
	filters := qps[1:]
	if len(filters) > 0 {
		filterStr = "@filter(" + strings.Join(filters, " AND ") + ")"
	}

	q = `
	{
		objects` + funcStr + filterStr + ` {
			uid
			expand(_all_) {
				uid
				expand(_all_)
			}
		}
	}
	`
	return q
}

// GetAllByClusterAndType - query to get result by filter edge
func (s DGClient) GetAllByClusterAndType(meta string, cluster string) (map[string]interface{}, error) {
	q := `
	{
  		objects (func: eq (objtype, "` + meta + `")) @cascade {
            uid
    		name
			resourceid
    		cluster @filter (eq(name, "` + cluster + `")) {
      			name
			}
  		}
	}`
	resp, err := s.dc.NewTxn().Query(context.Background(), q)
	if err != nil {
		log.Errorf("Query[%v] Error [%v]\n", q, err)
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(resp.Json, &m)
	if err != nil {
		log.Errorf("Query[%v] Error [%v]\n", q, err)
		return nil, err
	}
	return m, nil
}

// CreateSchema - create index
func (s DGClient) CreateSchema(sm Schema) error {
	var buffer bytes.Buffer
	buffer.WriteString(sm.Predicate)
	buffer.WriteString(": ")
	if sm.PType == "password" {
		buffer.WriteString(sm.PType)

	} else if sm.PType == util.UID {
		buffer.WriteString(sm.PType)
		if sm.Count {
			buffer.WriteString(" @count")
		}
		if sm.Reverse {
			buffer.WriteString(" @reverse")
		}
	} else {
		if sm.List {
			buffer.WriteString("[" + sm.PType + "]")
			if sm.Count {
				buffer.WriteString(" @count")
			}
		} else {
			buffer.WriteString(sm.PType)
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
