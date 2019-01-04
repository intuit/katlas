package apis

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/intuit/katlas/service/db"
	"github.com/intuit/katlas/service/util"
	"github.com/mitchellh/mapstructure"
)

// IMetaService define interfaces for metadata
type IMetaService interface {
	// Get metadata by name
	//GetMetadata(name string) (Metadata, error)
	GetMetadata(name string) (Metadata, error)
	// Create new metadata
	CreateMetadata(name string, data Metadata) (map[string]string, error)
	// Delete metadata
	DeleteMetadata(name string) error
	// Delete metadata field
	DeleteMetadataField(name string, filedName string) error
	// Append new field to metadata
	AddMetadataField(name string, data string) error
	// Get metadata field by name
	GetMetadataField(name string, fieldName string) (MetadataField, error)
	// Get all metadata fields
	GetMetadataFields(name string) ([]MetadataField, error)
}

// MetaService implements IMetaService interface
type MetaService struct {
	dbclient db.IDGClient
}

// Metadata Describe metadata
type Metadata struct {
	UID    string          `json:"uid"`
	Name   string          `json:"name"`
	Type   string          `json:"objtype,omitempty"`
	Fields []MetadataField `json:"fields,omitempty"`
	// TODO:
	// Add more needed attributes
}

// MetadataField describe attributes of Metadata
type MetadataField struct {
	FieldName string `json:"fieldName"`
	// Type of filed, could be one of [int, long, string, json, double, bool, date, enum, relationship]
	FieldType string `json:"fieldType"`
	// The field is required if value is true
	Mandatory bool `json:"mandatory"`
	// The field will be @index in schema if value is true
	Index bool `json:"index"`
	// If FieldType is relationship, need to set reference object type
	RefDataType string `json:"refDataType,omitempty"`
	// One or Many
	Cardinality string `json:"cardinality,omitempty"`
}

// GetMetadata get entity return the object with specified ID
func (s MetaService) GetMetadata(name string) (Metadata, error) {
	//var n Metadata
	qm := map[string][]string{util.Name: {name}, util.ObjType: {util.Metadata}}
	// Get metadata by name
	n := Metadata{}
	queryService := NewQueryService(s.dbclient)
	metas, err := queryService.GetQueryResult(qm)
	if err != nil {
		return n, err
	}
	if len(metas[util.Objects].([]interface{})) > 0 {
		meta := metas[util.Objects].([]interface{})[0].(map[string]interface{})
		var metadata Metadata
		err = mapstructure.Decode(meta, &metadata)
		if err != nil {
			return n, err
		}
		return metadata, nil
	}
	log.Debugf("GetMetadata return nil at end")
	return n, nil
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

func CheckKeys(keys []string, data map[string]interface{}) error {
	for k := range keys {
		if _, ok := data[keys[k]]; !ok {
			return fmt.Errorf("%q doesn't exist", keys[k])
		}
	}
	return nil
}

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
func (s MetaService) CreateMetadata(meta string, data map[string]interface{}) (map[string]string, error) {
	var rkeys = []string{"name", "fields", "objtype"}
	err := CheckKeys(rkeys, data)
	if err != nil {
		return nil, err
	}
	fMap := data["fields"].([]interface{})
	log.Debugf("fMap is %#v", fMap)
	const cardinality = "One"
	if len(fMap) > 0 {
		for i := range fMap {
			rkeys = []string{"fieldName", "fieldType"}
			err = CheckKeys(rkeys, fMap[i].(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			dkMap := map[string]interface{}{
				"cardinality": cardinality,
				"mandatory":   false,
				"index":       false,
			}
			err = SetDefaultKey(dkMap, fMap[i].(map[string]interface{}))
			if fMap[i].(map[string]interface{})["index"].(bool) == true {
				rkeys = []string{"upsert", "tokenizer"}
				err = CheckKeys(rkeys, fMap[i].(map[string]interface{}))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	e := NewEntityService(s.dbclient)
	uids, err := e.CreateEntity("metadata", data)

	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("can't create metadata")
	}
	log.Infof("metadata created/updated: %v", uids)
	return uids, nil
}
