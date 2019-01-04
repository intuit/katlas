package apis

import (
	"encoding/json"
	"fmt"
	"github.com/intuit/katlas/service/db"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestCreateEntity(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	//	ms := NewMetaService(dc)
	s := NewEntityService(dc)
	// create node
	node := map[string]interface{}{
		"objtype": "K8sNode",
		"name":    "node02",
		"labels":  "testingnode02",
	}
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "name", PType: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", PType: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", PType: "string", Index: true, Tokenizer: []string{"term"}})
	nids, _ := s.CreateEntity("K8sNode", node)
	var nid string
	for _, v := range nids {
		nid = v
		break
	}
	defer s.DeleteEntity(nid)
	n, _ := s.GetEntity("K8sNode", nid)
	o := n["objects"].([]interface{})[0].(map[string]interface{})
	if val, ok := o["labels"]; ok {
		assert.Equal(t, val, "testingnode02", "node label not equals to testnode02")
	} else {
		assert.Fail(t, "failed to create and get k8s node object")
	}
}

func TestDeleteEntityByRid(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	//	ms := NewMetaService(dc)
	s := NewEntityService(dc)
	// create node
	node := map[string]interface{}{
		"objtype":    "K8sNode",
		"name":       "node02",
		"labels":     "testingnode02",
		"resourceid": "noderid",
	}
	s.CreateEntity("K8sNode", node)
	err := s.DeleteEntityByResourceID("K8sNode", "noderid")
	assert.Nil(t, err)
	dc.Close()
}

func TestCreateEntityWithMeta(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	ms := NewMetaService(dc)
	s := NewEntityService(dc)
	q := NewQueryService(dc)
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "name", PType: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", PType: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", PType: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "namespace", PType: "uid", Tokenizer: []string{"term"}})

	podMeta := `{
		"name": "Pod",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldName": "name",
				"fieldType": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			},
			{
				"fieldName": "labels",
				"fieldType": "json",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			},
			{
				"fieldName": "objtype",
				"fieldType": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			},
			{
				"fieldName": "namespace",
				"fieldType": "relationship",
				"refDataType": "Namespace",
				"mandatory": true,
				"index": false,
				"cardinality": "One"
			},
			{
				"fieldName": "cluster",
				"fieldType": "relationship",
				"refDataType": "Cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "One"
			}
		]
	}`
	// create pod metadata
	dataMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(podMeta), &dataMap)
	if err != nil {
		panic(err)
	}
	s.CreateEntity("metadata", dataMap)
	// create namespace meta
	nsMeta := `{
		"name": "Namespace",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldName": "name",
				"fieldType": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			},
			{
				"fieldName": "cluster",
				"fieldType": "relationship",
				"refDataType": "Cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "One"
			}
		]
	}`
	nsMap := make(map[string]interface{})
	json.Unmarshal([]byte(nsMeta), &nsMap)
	s.CreateEntity("metadata", nsMap)
	clusterMeta := `{
		"name": "Cluster",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldName": "name",
				"fieldType": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			}
			
		]
	}`
	clMap := make(map[string]interface{})
	json.Unmarshal([]byte(clusterMeta), &clMap)
	s.CreateEntity("metadata", clMap)

	list := []string{"Pod", "Cluster", "Namespace"}
	for _, n := range list {
		// query to get created pod metadata
		qm := map[string][]string{"name": {n}, "objtype": {"metadata"}}
		n, _ := q.GetQueryResult(qm)
		o := n["objects"].([]interface{})[0].(map[string]interface{})
		// cleanup after test
		defer s.DeleteEntity(o["uid"].(string))
		for _, fields := range o["fields"].([]interface{}) {
			rid := fields.(map[string]interface{})["uid"]
			defer s.DeleteEntity(rid.(string))
		}
	}

	// create Pod data
	pod := `{
		"name": "pod01",
        "labels": {
			"app": "testpod",
            "label": "test"
		},
		"k8sobj": "K8sObj",
        "objtype": "Pod",
        "namespace": "default",
        "cluster": "cluster01",
		"resourceversion": "131"
	}`
	podMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(pod), &podMap)
	if err != nil {
		panic(err)
	}
	uids, err := s.CreateEntity("Pod", podMap)
	if err != nil {
		panic(err)
	}
	podMap["resourceversion"] = "132"
	podMap["namespace"] = "default"
	podMap["cluster"] = "cluster01"
	s.CreateEntity("Pod", podMap)

	for _, uid := range uids {
		pod, err := s.GetEntity("Pod", uid)
		if err != nil {
			assert.Fail(t, "Failed to get created Pod")
		}
		fs, _ := ms.GetMetadataFields("Pod")
		for _, f := range fs {
			if f.FieldType == "relationship" {
				for _, o := range pod["objects"].([]interface{}) {
					for _, r := range o.(map[string]interface{})[f.FieldName].([]interface{}) {
						s.dbclient.DeleteEntity(r.(map[string]interface{})["uid"].(string))
					}
				}
			}
		}
		s.dbclient.DeleteEntity(uid)
	}
}

func TestSyncEntities(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	ms := NewMetaService(dc)
	s := NewEntityService(dc)
	q := NewQueryService(dc)
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "name", PType: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", PType: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", PType: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "namespace", PType: "uid"})
	dc.CreateSchema(db.Schema{Predicate: "cluster", PType: "uid"})

	podMeta := `{
		"name": "Pod",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldName": "name",
				"fieldType": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			},
			{
				"fieldName": "labels",
				"fieldType": "json",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			},
			{
				"fieldName": "objtype",
				"fieldType": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			},
			{
				"fieldName": "namespace",
				"fieldType": "relationship",
				"refDataType": "Namespace",
				"mandatory": true,
				"index": false,
				"cardinality": "One"
			},
			{
				"fieldName": "cluster",
				"fieldType": "relationship",
				"refDataType": "Cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "One"
			}
		]
	}`
	// create pod metadata
	dataMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(podMeta), &dataMap)
	if err != nil {
		panic(err)
	}
	s.CreateEntity("metadata", dataMap)
	// create namespace meta
	nsMeta := `{
		"name": "Namespace",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldName": "name",
				"fieldType": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			},
			{
				"fieldName": "cluster",
				"fieldType": "relationship",
				"refDataType": "Cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "One"
			}
		]
	}`
	nsMap := make(map[string]interface{})
	json.Unmarshal([]byte(nsMeta), &nsMap)
	s.CreateEntity("metadata", nsMap)
	clusterMeta := `{
		"name": "Cluster",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldName": "name",
				"fieldType": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "One"
			}
			
		]
	}`
	clMap := make(map[string]interface{})
	json.Unmarshal([]byte(clusterMeta), &clMap)
	s.CreateEntity("metadata", clMap)

	list := []string{"Pod", "Cluster", "Namespace"}
	for _, n := range list {
		// query to get created pod metadata
		qm := map[string][]string{"name": {n}, "objtype": {"metadata"}}
		n, _ := q.GetQueryResult(qm)
		o := n["objects"].([]interface{})[0].(map[string]interface{})
		// cleanup after test
		defer s.DeleteEntity(o["uid"].(string))
		for _, fields := range o["fields"].([]interface{}) {
			rid := fields.(map[string]interface{})["uid"]
			defer s.DeleteEntity(rid.(string))
		}
	}

	// create Pod data
	pod := `{
		"name": "pod01",
        "labels": {
			"app": "testpod",
            "label": "test"
		},
		"k8sobj": "K8sObj",
        "objtype": "Pod",
        "namespace": "default01",
        "cluster": "cluster01",
		"resourceversion": "131"
	}`
	podMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(pod), &podMap)
	if err != nil {
		panic(err)
	}
	uids, err := s.CreateEntity("Pod", podMap)
	if err != nil {
		panic(err)
	}
	podMap["cluster"] = "cluster01"
	podMap["namespace"] = "default01"
	s.SyncEntities("Pod", []map[string]interface{}{podMap})
	for _, uid := range uids {
		pod, _ := s.GetEntity("Pod", uid)
		o := pod["objects"].([]interface{})[0].(map[string]interface{})
		assert.Equal(t, "131", o["resourceversion"].(string), "Pod got unexpected update")
	}
	// simulate sync pod with new version
	podMap["resourceversion"] = "132"
	podMap["cluster"] = "cluster01"
	podMap["namespace"] = "default01"
	s.SyncEntities("Pod", []map[string]interface{}{podMap})
	for _, uid := range uids {
		pod, _ := s.GetEntity("Pod", uid)
		o := pod["objects"].([]interface{})[0].(map[string]interface{})
		assert.Equal(t, "132", o["resourceversion"].(string), "Pod got unexpected update")
	}
	// simulate pod is not exist in k8s cluster and need to remove from database
	podMap["cluster"] = "cluster01"
	podMap["namespace"] = "default01"
	podMap["name"] = "pod02"
	s.SyncEntities("Pod", []map[string]interface{}{podMap})
	qm := map[string][]string{"name": {"pod01"}, "objtype": {"Pod"}}
	n, _ := q.GetQueryResult(qm)
	assert.Equal(t, 0, len(n["objects"].([]interface{})), "Pod01 still exist")
	qm2 := map[string][]string{"name": {"pod02"}, "objtype": {"Pod"}}
	n2, _ := q.GetQueryResult(qm2)
	o := n2["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "pod02", o["name"].(string), "Pod got unexpected creation")

	// sync namespace
	nsData := map[string]interface{}{
		"name":            "default01",
		"cluster":         "cluster01",
		"k8sobj":          "K8sObj",
		"objtype":         "Namespace",
		"resourceversion": "0",
	}
	s.SyncEntities("Namespace", []map[string]interface{}{nsData})
	qm3 := map[string][]string{"name": {"default01"}, "objtype": {"Namespace"}}
	n3, _ := q.GetQueryResult(qm3)
	o3 := n3["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "default01", o3["name"].(string), "Namespace default01 updated")

	nsData["cluster"] = "cluster01"
	nsData["name"] = "default02"
	s.SyncEntities("Namespace", []map[string]interface{}{nsData})
	qm4 := map[string][]string{"name": {"default01"}, "objtype": {"Namespace"}}
	n4, _ := q.GetQueryResult(qm4)
	assert.Equal(t, 0, len(n4["objects"].([]interface{})), "default namespace still exist")
	qm5 := map[string][]string{"name": {"default02"}, "objtype": {"Namespace"}}
	n5, _ := q.GetQueryResult(qm5)
	o5 := n5["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "default02", o5["name"].(string), "Namespace got unexpected creation")
	s.dbclient.DeleteEntity(o5["uid"].(string))

	for _, uid := range uids {
		pod, err := s.GetEntity("Pod", uid)
		if err != nil {
			assert.Fail(t, "Failed to get created Pod")
		}
		fs, _ := ms.GetMetadataFields("Pod")
		for _, f := range fs {
			if f.FieldType == "relationship" {
				for _, o := range pod["objects"].([]interface{}) {
					for _, r := range o.(map[string]interface{})[f.FieldName].([]interface{}) {
						s.dbclient.DeleteEntity(r.(map[string]interface{})["uid"].(string))
					}
				}
			}
		}
		s.dbclient.DeleteEntity(uid)
	}
}

func TestMultiCreateEntity(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	q := NewQueryService(dc)
	defer dc.Close()
	//	ms := NewMetaService(dc)
	s := NewEntityService(dc)

	var wg sync.WaitGroup
	rest := make(chan map[string]string)
	wg.Add(100)
	for i := 0; i < 50; i++ {
		go func(version string) {
			defer wg.Done()
			node := map[string]interface{}{
				"objtype":    "K8sNode",
				"name":       "multinode",
				"label":      version,
				"resourceid": "cluster:ns:multinode",
			}
			nids, _ := s.CreateEntity("K8sNode", node)
			rest <- nids
		}(strconv.Itoa(i))
	}

	for i := 0; i < 50; i++ {
		go func(version string) {
			defer wg.Done()
			node := map[string]interface{}{
				"objtype":         "K8sNode",
				"name":            "multinode2",
				"label":           version,
				"resourceid":      "cluster:ns:multinode2",
				"resourceversion": version,
			}
			nids, _ := s.CreateEntity("K8sNode", node)
			rest <- nids
		}(strconv.Itoa(i))
	}

	go func() {
		for r := range rest {
			fmt.Println(r)
		}
	}()
	wg.Wait()
	qm := map[string][]string{"resourceid": {"cluster:ns:multinode"}, "objtype": {"K8sNode"}}
	n, _ := q.GetQueryResult(qm)
	o := n["objects"].([]interface{})
	assert.Equal(t, 1, len(o), "only one object expect to be created with same resourceid")
	s.DeleteEntityByResourceID("K8sNode", "cluster:ns:multinode")
	s.DeleteEntityByResourceID("K8sNode", "cluster:ns:multinode2")
}
