package apis

import (
	"encoding/json"
	"testing"

	"github.com/intuit/katlas/service/db"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestMetaService(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")

	q := NewQueryService(dc)
	m := NewMetaService(dc)
	e := NewEntityService(dc)
	// create pod metadata
	podMeta := `{
		"name": "pod_metadata",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "json",
				"mandatory": true,
				"cardinality": "one"
			},
			{
				"fieldname": "status",
				"fieldtype": "string",
				"mandatory": true,
				"cardinality": "one"
			},
			{
				"fieldname": "containers",
				"fieldtype": "relationship",
				"refdatatype": "K8scontainer",
				"mandatory": false,
				"cardinality": "many"
			}
		]
	}`
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	// create pod metadata
	dataMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(podMeta), &dataMap)
	if err != nil {
		panic(err)
	}
	m.CreateMetadata(dataMap)
	// query to get created pod metadata
	qm := map[string][]string{"name": {"pod_metadata"}, "objtype": {"metadata"}}
	n, _ := q.GetQueryResult(qm)
	o := n["objects"].([]interface{})[0].(map[string]interface{})
	// cleanup after test
	defer e.DeleteEntity(o["uid"].(string))
	assert.Equal(t, o["name"], "pod_metadata", "query return doesn't match pod_metadata")
	for _, fields := range o["fields"].([]interface{}) {
		rid := fields.(map[string]interface{})["uid"]
		defer e.DeleteEntity(rid.(string))
	}
	// get all fields
	fs, err := m.GetMetadataFields("pod_metadata")
	if err != nil {
		assert.Fail(t, "Failed to get meta fields")
	}
	assert.Equal(t, 3, len(fs), "return fields don't match metadata define")
}

func TestDeleteMetadata(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	m := NewMetaService(dc)
	q := NewQueryService(dc)
	// create pod metadata
	podMeta := `{
		"name": "pod_metadata",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "json",
				"mandatory": true,
				"cardinality": "one"
			},
			{
				"fieldname": "status",
				"fieldtype": "string",
				"mandatory": true,
				"cardinality": "one"
			},
			{
				"fieldname": "containers",
				"fieldtype": "relationship",
				"refdatatype": "k8container",
				"mandatory": false,
				"cardinality": "many"
			}
		]
	}`
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	// create pod metadata
	dataMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(podMeta), &dataMap)
	if err != nil {
		panic(err)
	}
	m.CreateMetadata(dataMap)
	// create referenced metadata
	k8container := `{
		"name": "k8container",
		"objtype" : "metadata",
		"fields": [
		{
			"fieldname": "name",
			"fieldtype": "json",
			"mandatory": true,
			"cardinality": "one"
		}
	]
	}`
	cmap := make(map[string]interface{})
	err = json.Unmarshal([]byte(k8container), &cmap)
	if err != nil {
		panic(err)
	}
	m.CreateMetadata(cmap)
	// delete fail due to meta been referenced
	error := m.DeleteMetadata("k8container")
	assert.NotNil(t, error)
	// remove pod metadata
	error = m.DeleteMetadata("pod_metadata")
	qm := map[string][]string{"name": {"pod_metadata"}, "objtype": {"metadata"}}
	n, _ := q.GetQueryResult(qm)
	assert.Empty(t, n["objects"])
	// clean
	m.DeleteMetadata("k8container")
}

func TestMetadataUpdate(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	m := NewMetaService(dc)
	q := NewQueryService(dc)
	// create pod metadata
	podMeta := `{
		"name": "pod_metadata",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "json",
				"mandatory": true,
				"cardinality": "one"
			},
			{
				"fieldname": "status",
				"fieldtype": "string",
				"mandatory": true,
				"cardinality": "one"
			}
		]
	}`
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	// create pod metadata
	dataMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(podMeta), &dataMap)
	if err != nil {
		panic(err)
	}
	m.CreateMetadata(dataMap)
	fs := make([]interface{}, 0)
	f := map[string]interface{}{
		"fieldname": "name",
		"fieldtype": "string",
	}
	fs = append(fs, f)
	m.UpdateMetadata("pod_metadata", map[string]interface{}{
		"fields": fs,
	})
	qm := map[string][]string{"name": {"pod_metadata"}, "objtype": {"metadata"}}
	n, _ := q.GetQueryResult(qm)
	o := n["objects"].([]interface{})[0].(map[string]interface{})
	var metadata Metadata
	mapstructure.Decode(o, &metadata)
	for _, f := range metadata.Fields {
		assert.Equal(t, "string", f.FieldType, "fieldtype should updated to string")
	}
	m.DeleteMetadata("pod_metadata")
}
