package apis

import (
	"encoding/json"
	"strings"
	"time"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/cenkalti/backoff"
	"github.com/intuit/katlas/service/db"
	"github.com/intuit/katlas/service/metrics"
	"github.com/intuit/katlas/service/util"
	"reflect"
)

// KeyMutex lock single object by resourceid
// It will retry to get lock and expire after certain time
var mutex = util.NewKeyMutex(time.Minute, 1000)

// IEntityService define interfaces to manipulate data
type IEntityService interface {
	// get entity return the object with specified ID
	GetEntity(uid string) (map[string]interface{}, error)
	// remove object with given ID
	DeleteEntity(uid string) error
	// remove object with given ID
	DeleteEntityByResourceID(rid string) error
	// save new entity to the storage
	CreateEntity(meta string, data map[string]interface{}) (string, error)
	// update entity with given ID in the storage
	UpdateEntity(uuid string, data map[string]interface{}, option ...util.OptionContext)
	// create or remove relationship between entities by given IDs
	CreateOrDeleteEdge(fromUID string, toUID string, rel string, op db.Action) error
	// sync data between source and underlying database
	SyncEntities(meta string, data map[string]interface{}) error
}

// EntityService provides service for controller and frontend by implement IEntityService interface
type EntityService struct {
	dbclient db.IDGClient
}

// NewEntityService creates a new EntityService with the given dgraph client.
func NewEntityService(dc db.IDGClient) *EntityService {
	return &EntityService{dc}
}

// GetEntity get entity return the object with specified ID
func (s EntityService) GetEntity(uuid string) (map[string]interface{}, error) {
	metrics.DgraphNumGetEntity.Inc()
	return s.dbclient.GetEntity(uuid)
}

// DeleteEntity remove object with given ID
func (s EntityService) DeleteEntity(uuid string) error {
	metrics.DgraphNumDeleteEntity.Inc()
	return s.dbclient.DeleteEntity(uuid)
}

// DeleteEntityByResourceID remove object by given resourceid
func (s EntityService) DeleteEntityByResourceID(meta string, rid string) error {
	metrics.DgraphNumDeleteEntity.Inc()
	qm := map[string][]string{util.ResourceID: {rid}, util.ObjType: {meta}, util.Print: {util.ResourceID}}
	queryService := NewQueryService(s.dbclient)
	node, err := queryService.GetQueryResult(qm)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(node[util.Objects].([]interface{})) > 0 {
		// got existing object id
		for _, obj := range node[util.Objects].([]interface{}) {
			err = s.dbclient.DeleteEntity(obj.(map[string]interface{})[util.UID].(string))
			if err != nil {
				return err
			}
			log.Infof("%s %s deleted successfully", meta, rid)
		}
	}
	return nil
}

// CreateEntity save new entity to the storage
func (s EntityService) CreateEntity(meta string, data map[string]interface{}) (string, error) {
	metrics.DgraphNumCreateEntity.Inc()
	cluster := data[util.Cluster]
	ns := data[util.Namespace]
	if _, ok := data[util.ResourceID]; !ok {
		data[util.ResourceID] = getResourceID(meta, data)
	}

	m := NewMetaService(s.dbclient)
	fs, err := m.GetMetadataFields(meta)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	delMap := make(map[string]interface{})
	if len(fs) > 0 {
		for _, field := range fs {
			fieldValue, ok := data[field.FieldName]
			if !ok || fieldValue == "" || fieldValue == nil || ((reflect.ValueOf(fieldValue).Kind() == reflect.Interface ||
				reflect.ValueOf(fieldValue).Kind() == reflect.Ptr ||
				reflect.ValueOf(fieldValue).Kind() == reflect.Slice) &&
				reflect.ValueOf(fieldValue).IsNil()) {
				delete(data, field.FieldName)
				continue
			}
			if strings.EqualFold(field.Cardinality, util.Many) {
				delMap[field.FieldName] = nil
			}
			// if field type is json, should convert it to facets - key value pair on node
			// simply convert it to string for now
			if field.FieldType == util.JSON {
				bytes, _ := json.Marshal(data[field.FieldName])
				data[field.FieldName] = string(bytes)
			} else if field.FieldType == util.Relationship {
				// replace value with relationship object uid
				if strings.EqualFold(field.Cardinality, util.Many) {
					uidMaps := []map[string]interface{}{}
					for _, rel := range data[field.FieldName].([]interface{}) {
						dataMap := buildDataMap(data[util.K8sObj], rel, field.RefDataType, cluster, ns)
						uid, err := s.getUIDFromRelData(dataMap, field.RefDataType)
						if err != nil {
							log.Error(err)
							return "", err
						}
						uidMaps = append(uidMaps, map[string]interface{}{util.UID: *uid})
					}
					data[field.FieldName] = uidMaps
				} else {
					uidMap := map[string]interface{}{}
					// pod can be owned by multi-objs like replicaset, daemonset
					// FIX-later: hack to set refDataType to dynamic value from owner reference
					if strings.EqualFold(meta, util.Pod) && strings.EqualFold(field.FieldName, util.Owner) {
						field.RefDataType = data[util.OwnerType].(string)
					}
					dataMap := buildDataMap(data[util.K8sObj], data[field.FieldName], field.RefDataType, cluster, ns)
					uid, err := s.getUIDFromRelData(dataMap, field.RefDataType)
					if err != nil {
						log.Error(err)
						return "", err
					}
					uidMap[util.UID] = *uid
					data[field.FieldName] = uidMap
				}
			}
		}
	}
	if _, ok := data[util.UID]; !ok {
		data[util.UID] = "_:A"
	}
	if mutex.TryLock(data[util.ResourceID]) {
		defer mutex.Unlock(data[util.ResourceID])
		var uuid string
		operation := func() error {
			uuid, err = s.dbclient.CreateEntity(meta, data)
			return err
		}
		err := backoff.Retry(operation, backoff.WithMaxRetries(util.NewBackOff(), util.RetryCount))
		if err != nil {
			return "", err
		}
		return uuid, nil
	}
	return "", fmt.Errorf("can't get resource lock, ignore after timeout reached")
}

// SyncEntities ...
func (s EntityService) SyncEntities(meta string, data []map[string]interface{}) error {
	// get all objects from database base on meta and k8 cluster
	if len(data) > 0 {
		objs, err := s.dbclient.GetAllByClusterAndType(meta, data[0][util.Cluster].(string))
		if err != nil {
			log.Error(err)
			return err
		}
		// if obj not present in input, remove it from database
		if len(objs[util.Objects].([]interface{})) > 0 {
			for _, obj := range objs[util.Objects].([]interface{}) {
				uid := obj.(map[string]interface{})[util.UID].(string)
				rid := obj.(map[string]interface{})[util.ResourceID].(string)
				name := obj.(map[string]interface{})[util.Name].(string)
				var ns interface{}
				if strings.Count(rid, ":") == 3 {
					ns = strings.Split(rid, ":")[2]
				}
				found := false
				for _, d := range data {
					if name == d[util.Name].(string) && ns == d[util.Namespace] {
						found = true
					}
				}
				if !found {
					s.dbclient.DeleteEntity(uid)
				}
			}
		}
		// create or update from input
		for _, d := range data {
			s.CreateEntity(meta, d)
		}
	}
	return nil
}

// CreateOrDeleteEdge create or remove edge
func (s EntityService) CreateOrDeleteEdge(fromType string, fromUID string, toType, toUID string, rel string, op db.Action) error {
	// TODO:
	// validate base on metadata
	// if err := metadata.Validate(fromType, toType, rel); err != nil {
	// 	return nil, err
	// }
	metrics.DgraphNumUpdateEdge.Inc()
	return s.dbclient.CreateOrDeleteEdge(fromType, fromUID, toType, toUID, rel, op)

}

// UpdateEntity update entity
func (s EntityService) UpdateEntity(uuid string, data map[string]interface{}, option ...util.OptionContext) error {
	if mutex.TryLock(uuid) {
		defer mutex.Unlock(uuid)
		operation := func() error {
			return s.dbclient.UpdateEntity(uuid, data, option...)
		}
		err := backoff.Retry(operation, backoff.WithMaxRetries(util.NewBackOff(), util.RetryCount))
		if err != nil {
			return err
		}
		metrics.DgraphNumUpdateEntity.Inc()
		return nil
	}
	return fmt.Errorf("can't get resource lock to update %s, ignore after timeout reached", uuid)
}

// build resourceid
func getResourceID(meta string, data map[string]interface{}) string {
	ridPrefix := meta + ":"
	if _, ok := data[util.K8sObj]; ok {
		if _, ok := data[util.Cluster]; ok {
			ridPrefix += data[util.Cluster].(string) + ":"
		}
		if _, ok := data[util.Namespace]; ok {
			ridPrefix += data[util.Namespace].(string) + ":"
		}
	}
	return ridPrefix + data[util.Name].(string)
}

// build data
func buildDataMap(k8sObj interface{}, relData interface{}, relType string, cluster interface{}, ns interface{}) map[string]interface{} {
	var dataMap map[string]interface{}
	if reflect.TypeOf(relData).Kind() == reflect.String {
		dataMap = make(map[string]interface{})
		dataMap[util.Name] = relData.(string)
	} else if reflect.TypeOf(relData).Kind() == reflect.Map {
		dataMap = relData.(map[string]interface{})
	}
	dataMap[util.ObjType] = relType
	if k8sObj != nil {
		dataMap[util.K8sObj] = util.K8sObj
	}
	_, hasRID := dataMap[util.ResourceID]
	_, hasUID := dataMap[util.UID]
	if hasRID || hasUID {
		return dataMap
	}
	// compose resource id
	if strings.EqualFold(relType, util.Cluster) {
		dataMap[util.ResourceID] = relType + ":" + dataMap[util.Name].(string)
	} else if strings.EqualFold(relType, util.Namespace) || strings.EqualFold(relType, util.Node) {
		dataMap[util.Cluster] = cluster
		dataMap[util.ResourceID] = relType + ":" + cluster.(string) + ":" + dataMap[util.Name].(string)
	} else if strings.EqualFold(relType, util.Application) || strings.EqualFold(relType, util.Asset) {
		dataMap[util.ResourceID] = relType + ":" + dataMap[util.Name].(string)
	} else {
		if cluster != nil {
			dataMap[util.Cluster] = cluster
		}
		if ns != nil {
			dataMap[util.Namespace] = ns
		}
		dataMap[util.ResourceID] = getResourceID(relType, dataMap)
	}
	return dataMap
}

// get uid from relationship object, if object not present, create it
func (s EntityService) getUIDFromRelData(data map[string]interface{}, objType string) (*string, error) {
	if _, ok := data[util.UID]; ok {
		id := data[util.UID].(string)
		return &id, nil
	}
	// query by ResourceID to get uid
	qm := map[string][]string{util.ResourceID: {data[util.ResourceID].(string)}, util.ObjType: {objType}, util.Print: {util.ResourceID}}
	queryService := NewQueryService(s.dbclient)
	node, err := queryService.GetQueryResult(qm)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var uid string
	if len(node[util.Objects].([]interface{})) > 0 {
		// got existing object id
		uid = node[util.Objects].([]interface{})[0].(map[string]interface{})[util.UID].(string)
	} else {
		// create new object
		uid, err = s.CreateEntity(objType, data)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return &uid, nil
}
