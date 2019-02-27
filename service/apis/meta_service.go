package apis

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/golang-lru"
	"github.com/intuit/katlas/service/db"
	"github.com/intuit/katlas/service/util"
	"github.com/mitchellh/mapstructure"
	"strings"
)

// IMetaService define interfaces for metadata
type IMetaService interface {
	// Get metadata by name
	GetMetadata(name string) (*Metadata, error)
	// Create new metadata
	CreateMetadata(data Metadata) (string, error)
	// Delete metadata
	DeleteMetadata(name string) error
	// Update metadata
	UpdateMetadata(name string, data map[string]interface{}) error
	// Get all metadata fields
	GetMetadataFields(name string) ([]MetadataField, error)
	// Create schema
	CreateSchema(sm db.Schema) error
	// Drop a schema
	DropSchema(name string) error
	// Remove schema from cache
	RemoveSchemaCache(cache *lru.Cache)
}

// MetaService implements IMetaService interface
type MetaService struct {
	dbclient db.IDGClient
}

// Metadata Describe metadata
type Metadata struct {
	UID             string          `json:"uid"`
	Name            string          `json:"name"`
	ObjType         string          `json:"objtype,omitempty"`
	Fields          []MetadataField `json:"fields,omitempty"`
	ResourceVersion string          `json:"resourceversion,omitempty"`
	// TODO:
	// Add more needed attributes
}

// MetadataField describe attributes of Metadata
type MetadataField struct {
	UID       string `json:"uid"`
	FieldName string `json:"fieldname"`
	// Type of filed, could be one of [int, long, string, json, double, bool, date, enum, relationship]
	FieldType string `json:"fieldtype"`
	// The field is required if value is true
	Mandatory bool `json:"mandatory"`
	// If FieldType is relationship, need to set reference object type
	RefDataType string `json:"refdatatype,omitempty"`
	// One or Many
	Cardinality string `json:"cardinality,omitempty"`
}

// GetMetadata get entity return the object with specified ID
func (s MetaService) GetMetadata(name string) (*Metadata, error) {
	//var n Metadata
	qm := map[string][]string{util.Name: {name}, util.ObjType: {util.Metadata}}
	// Get metadata by name
	queryService := NewQueryService(s.dbclient)
	metas, err := queryService.GetQueryResult(qm)
	if err != nil {
		return nil, err
	}
	if len(metas[util.Objects].([]interface{})) > 0 {
		meta := metas[util.Objects].([]interface{})[0].(map[string]interface{})
		var metadata Metadata
		err = mapstructure.Decode(meta, &metadata)
		if err != nil {
			return nil, err
		}
		return &metadata, nil
	}
	return nil, nil
}

// NewMetaService creates a new MetaService with the given dgraph client.
func NewMetaService(dc db.IDGClient) *MetaService {
	return &MetaService{dc}
}

// GetMetadataFields EntityService will call this method to get all fields to verify and create edge accordingly
func (s MetaService) GetMetadataFields(name string) ([]MetadataField, error) {
	qm := map[string][]string{util.Name: {name}, util.ObjType: {util.Metadata}}
	// Get metadata by name
	queryService := NewQueryService(s.dbclient)
	metas, err := queryService.GetQueryResult(qm)
	if err != nil {
		return nil, err
	}
	if len(metas[util.Objects].([]interface{})) > 0 {
		meta := metas[util.Objects].([]interface{})[0].(map[string]interface{})
		var metadata Metadata
		err = mapstructure.Decode(meta, &metadata)
		if err != nil {
			return nil, err
		}
		return metadata.Fields, nil
	}
	return nil, nil
}

// CheckKeys checks if keys exist
func CheckKeys(keys []string, data map[string]interface{}) error {
	for k := range keys {
		if _, ok := data[keys[k]]; !ok {
			return fmt.Errorf("%q doesn't exist", keys[k])
		}
	}
	return nil
}

//SetDefaultKey sets default values for the keys if any isn't set
func SetDefaultKey(dkMap map[string]interface{}, data map[string]interface{}) error {
	for key := range dkMap {
		if _, ok := data[key]; !ok {
			data[key] = dkMap[key]
			log.Infof("%v doesn't exist. set to default %#v", key, dkMap[key])
		}
	}
	return nil
}

// CreateMetadata save new metadata to the storage
func (s MetaService) CreateMetadata(data map[string]interface{}) (string, error) {
	queryService := NewQueryService(s.dbclient)
	qm := map[string][]string{util.Name: {data[util.Name].(string)}, util.ObjType: {util.Metadata}}
	metas, _ := queryService.GetQueryResult(qm)
	if len(metas[util.Objects].([]interface{})) > 0 {
		return "", fmt.Errorf("metadata %s already exist, creation failed", data[util.Name].(string))
	}
	var rkeys = []string{util.Name, util.Fields, util.ObjType}
	err := CheckKeys(rkeys, data)
	if err != nil {
		return "", err
	}
	fMap, ok := data[util.Fields].([]interface{})
	if !ok {
		return "", fmt.Errorf("error in metadata field")
	}

	if len(fMap) > 0 {
		for i := range fMap {
			rkeys = []string{util.FieldName, util.FieldType}
			err = CheckKeys(rkeys, fMap[i].(map[string]interface{}))
			if err != nil {
				return "", err
			}
			dkMap := map[string]interface{}{
				util.Cardinality: util.One,
				util.Mandatory:   false,
			}
			err = SetDefaultKey(dkMap, fMap[i].(map[string]interface{}))
		}
	}
	e := NewEntityService(s.dbclient)
	uid, err := e.CreateEntity(util.Metadata, data)

	if err != nil {
		log.Error(err)
		return "", fmt.Errorf("can't create metadata %v", err)
	}
	return uid, nil
}

// CreateSchema creates schema
func (s MetaService) CreateSchema(sm db.Schema) error {
	return s.dbclient.CreateSchema(sm)
}

// DropSchema to remove schema
func (s MetaService) DropSchema(name string) error {
	return s.dbclient.DropSchema(name)
}

// RemoveSchemaCache to clean lru cache
func (s MetaService) RemoveSchemaCache(cache *lru.Cache) {
	s.dbclient.RemoveDBSchemaFromCache(cache)
}

// DeleteMetadata to remove metadata if not been referenced by others
func (s MetaService) DeleteMetadata(name string) error {
	qm := map[string][]string{util.ObjType: {util.Metadata}}
	// Get metadata by name
	queryService := NewQueryService(s.dbclient)
	metas, err := queryService.GetQueryResult(qm)
	if err != nil {
		return err
	}
	var target Metadata
	for _, meta := range metas[util.Objects].([]interface{}) {
		var metadata Metadata
		err := mapstructure.Decode(meta, &metadata)
		if err != nil {
			return err
		}
		for _, field := range metadata.Fields {
			if strings.Contains(field.RefDataType, name) {
				return fmt.Errorf("not able to delete metadata %s which is referenced by %s", name, metadata.Name)
			}
		}
		if metadata.Name == name {
			target = metadata
		}
	}
	for _, field := range target.Fields {
		err = s.dbclient.DeleteEntity(field.UID)
		if err != nil {
			return err
		}
	}
	return s.dbclient.DeleteEntity(target.UID)
}

// UpdateMetadata update metadata fields
// The payload can include either new fields or existing fields
func (s MetaService) UpdateMetadata(name string, data map[string]interface{}) error {
	qm := map[string][]string{util.Name: {name}, util.ObjType: {util.Metadata}}
	queryService := NewQueryService(s.dbclient)
	metas, err := queryService.GetQueryResult(qm)
	if err != nil {
		return err
	}
	if len(metas[util.Objects].([]interface{})) > 0 {
		meta := metas[util.Objects].([]interface{})[0].(map[string]interface{})
		var metadata Metadata
		err = mapstructure.Decode(meta, &metadata)
		if err != nil {
			return err
		}
		data[util.UID] = metadata.UID
		if fieldMap, ok := data[util.Fields]; ok {
			for _, fs := range fieldMap.([]interface{}) {
				name := fs.(map[string]interface{})[util.FieldName]
				for _, currentField := range metadata.Fields {
					if currentField.FieldName == name {
						// set field uid with existing one
						// this will enable single metadata filed can be updated
						fs.(map[string]interface{})[util.UID] = currentField.UID
						break
					}
				}
			}
		}
		e := NewEntityService(s.dbclient)
		err := e.UpdateEntity(metadata.UID, data, util.OptionContext{ReplaceListOrEdge: false})
		return err
	}
	return fmt.Errorf("metadata %s not found", name)
}
