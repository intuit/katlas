package apis

import (
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
	metas, err := s.dbclient.GetQueryResult(qm)
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
func NewMetaService(dc *db.DGClient) *MetaService {
	return &MetaService{dc}
}

// GetMetadataFields EntityService will call this method to get all fields to verify and create edge accordingly
func (s MetaService) GetMetadataFields(name string) ([]MetadataField, error) {
	qm := map[string][]string{util.Name: {name}, util.ObjType: {util.Metadata}}
	// Get metadata by name
	metas, err := s.dbclient.GetQueryResult(qm)
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
