package apis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/intuit/katlas/service/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateEntity(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	// create node
	node := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "node02",
		"labels":  "testingnode02",
	}
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	nid, _ := s.CreateEntity("k8snode", node)
	defer s.DeleteEntity(nid)
	n, _ := s.GetEntity(nid)
	o := n["objects"].([]interface{})[0].(map[string]interface{})
	if val, ok := o["labels"]; ok {
		assert.Equal(t, val, "testingnode02", "node label not equals to testnode02")
	} else {
		assert.Fail(t, "failed to create and get k8s node object")
	}
}

func TestDeleteEntityByRid(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	s := NewEntityService(dc)
	// create node
	node := map[string]interface{}{
		"objtype":    "k8snode",
		"name":       "node02",
		"labels":     "testingnode02",
		"resourceid": "noderid",
	}
	s.CreateEntity("k8snode", node)
	err := s.DeleteEntityByResourceID("k8snode", "noderid")
	assert.Nil(t, err)
	dc.Close()
}

func TestCreateEntityWithMeta(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	q := NewQueryService(dc)
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "namespace", Type: "uid"})
	dc.CreateSchema(db.Schema{Predicate: "cluster", Type: "uid"})

	podMeta := `{
		"name": "pod",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "labels",
				"fieldtype": "json",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "objtype",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "namespace",
				"fieldtype": "relationship",
				"refdatatype": "namespace",
				"mandatory": true,
				"index": false,
				"cardinality": "one"
			},
			{
				"fieldname": "cluster",
				"fieldtype": "relationship",
				"refdatatype": "cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "one"
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
		"name": "namespace",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "cluster",
				"fieldtype": "relationship",
				"refdatatype": "cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "one"
			}
		]
	}`
	nsMap := make(map[string]interface{})
	json.Unmarshal([]byte(nsMeta), &nsMap)
	s.CreateEntity("metadata", nsMap)
	clusterMeta := `{
		"name": "cluster",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			}
			
		]
	}`
	clMap := make(map[string]interface{})
	json.Unmarshal([]byte(clusterMeta), &clMap)
	s.CreateEntity("metadata", clMap)

	list := []string{"pod", "cluster", "namespace"}
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

	// create pod data
	pod := `{
		"name": "pod01",
        "labels": {
			"app": "testpod",
            "label": "test"
		},
		"k8sobj": "K8sObj",
        "objtype": "pod",
        "namespace": "default",
        "cluster": "cluster01",
		"resourceversion": "131"
	}`
	podMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(pod), &podMap)
	if err != nil {
		panic(err)
	}
	uid, err := s.CreateEntity("pod", podMap)
	if err != nil {
		panic(err)
	}
	podMap["resourceversion"] = "132"
	podMap["namespace"] = "default"
	podMap["cluster"] = "cluster01"
	s.CreateEntity("pod", podMap)

	s.DeleteEntityByResourceID("namespace", "namespace:cluster01:default")
	s.DeleteEntityByResourceID("cluster", "cluster:cluster01")
	s.dbclient.DeleteEntity(uid)
}

func TestSyncEntities(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	ms := NewMetaService(dc)
	s := NewEntityService(dc)
	q := NewQueryService(dc)
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "namespace", Type: "uid"})
	dc.CreateSchema(db.Schema{Predicate: "cluster", Type: "uid"})

	podMeta := `{
		"name": "pod",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "labels",
				"fieldtype": "json",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "objtype",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "namespace",
				"fieldtype": "relationship",
				"refdatatype": "namespace",
				"mandatory": true,
				"index": false,
				"cardinality": "one"
			},
			{
				"fieldname": "cluster",
				"fieldtype": "relationship",
				"refdatatype": "cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "one"
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
		"name": "namespace",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "cluster",
				"fieldtype": "relationship",
				"refdatatype": "cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "one"
			}
		]
	}`
	nsMap := make(map[string]interface{})
	json.Unmarshal([]byte(nsMeta), &nsMap)
	s.CreateEntity("metadata", nsMap)
	clusterMeta := `{
		"name": "cluster",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			}
			
		]
	}`
	clMap := make(map[string]interface{})
	json.Unmarshal([]byte(clusterMeta), &clMap)
	s.CreateEntity("metadata", clMap)

	list := []string{"pod", "cluster", "namespace"}
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

	// create pod data
	pod := `{
		"name": "pod01",
        "labels": {
			"app": "testpod",
            "label": "test"
		},
		"k8sobj": "K8sObj",
        "objtype": "pod",
        "namespace": "default01",
        "cluster": "cluster01",
		"resourceversion": "131"
	}`
	podMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(pod), &podMap)
	if err != nil {
		panic(err)
	}
	uid, err := s.CreateEntity("pod", podMap)
	if err != nil {
		panic(err)
	}
	podMap["cluster"] = "cluster01"
	podMap["namespace"] = "default01"
	s.SyncEntities("pod", []map[string]interface{}{podMap})

	pods, _ := s.GetEntity(uid)
	o := pods["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "131", o["resourceversion"].(string), "pod got unexpected update")

	// simulate sync pod with new version
	podMap["resourceversion"] = "132"
	podMap["cluster"] = "cluster01"
	podMap["namespace"] = "default01"
	s.SyncEntities("pod", []map[string]interface{}{podMap})

	pod2, _ := s.GetEntity(uid)
	o = pod2["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "132", o["resourceversion"].(string), "pod got unexpected update")

	// simulate pod is not exist in k8s cluster and need to remove from database
	podMap["cluster"] = "cluster01"
	podMap["namespace"] = "default01"
	podMap["name"] = "pod02"
	s.SyncEntities("pod", []map[string]interface{}{podMap})
	qm := map[string][]string{"name": {"pod01"}, "objtype": {"pod"}}
	n, _ := q.GetQueryResult(qm)
	assert.Equal(t, 0, len(n["objects"].([]interface{})), "pod01 still exist")
	qm2 := map[string][]string{"name": {"pod02"}, "objtype": {"pod"}}
	n2, _ := q.GetQueryResult(qm2)
	o = n2["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "pod02", o["name"].(string), "pod got unexpected creation")

	// sync namespace
	nsData := map[string]interface{}{
		"name":            "default01",
		"cluster":         "cluster01",
		"k8sobj":          "K8sObj",
		"objtype":         "namespace",
		"resourceversion": "0",
	}
	s.SyncEntities("namespace", []map[string]interface{}{nsData})
	qm3 := map[string][]string{"name": {"default01"}, "objtype": {"namespace"}}
	n3, _ := q.GetQueryResult(qm3)
	o3 := n3["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "default01", o3["name"].(string), "namespace default01 updated")

	nsData["cluster"] = "cluster01"
	nsData["name"] = "default02"
	s.SyncEntities("namespace", []map[string]interface{}{nsData})
	qm4 := map[string][]string{"name": {"default01"}, "objtype": {"namespace"}}
	n4, _ := q.GetQueryResult(qm4)
	assert.Equal(t, 0, len(n4["objects"].([]interface{})), "default namespace still exist")
	qm5 := map[string][]string{"name": {"default02"}, "objtype": {"namespace"}}
	n5, _ := q.GetQueryResult(qm5)
	o5 := n5["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "default02", o5["name"].(string), "namespace got unexpected creation")
	s.dbclient.DeleteEntity(o5["uid"].(string))

	pod3, err := s.GetEntity(uid)
	if err != nil {
		assert.Fail(t, "Failed to get created pod")
	}
	fs, _ := ms.GetMetadataFields("pod")
	for _, f := range fs {
		if f.FieldType == "relationship" {
			for _, o := range pod3["objects"].([]interface{}) {
				for _, r := range o.(map[string]interface{})[f.FieldName].([]interface{}) {
					s.dbclient.DeleteEntity(r.(map[string]interface{})["uid"].(string))
				}
			}
		}
	}
	s.dbclient.DeleteEntity(uid)

}

func TestMultiCreateEntity(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	q := NewQueryService(dc)
	defer dc.Close()
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Upsert: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	s := NewEntityService(dc)
	var wg sync.WaitGroup
	rest := make(chan string)
	wg.Add(100)
	for i := 0; i < 50; i++ {
		go func(version string) {
			defer wg.Done()
			node := map[string]interface{}{
				"objtype":    "k8snode",
				"name":       "multinode",
				"label":      version,
				"resourceid": "cluster:ns:multinode",
			}
			nid, _ := s.CreateEntity("k8snode", node)
			rest <- nid
		}(strconv.Itoa(i))
	}

	for i := 0; i < 50; i++ {
		go func(version string) {
			defer wg.Done()
			node := map[string]interface{}{
				"objtype":         "k8snode",
				"name":            "multinode2",
				"label":           version,
				"resourceid":      "cluster:ns:multinode2",
				"resourceversion": version,
			}
			nid, _ := s.CreateEntity("k8snode", node)
			rest <- nid
		}(strconv.Itoa(i))
	}

	go func() {
		for r := range rest {
			fmt.Println(r)
		}
	}()
	wg.Wait()
	qm := map[string][]string{"resourceid": {"cluster:ns:multinode"}, "objtype": {"k8snode"}}
	n, _ := q.GetQueryResult(qm)
	o := n["objects"].([]interface{})
	assert.Equal(t, 1, len(o), "only one object expect to be created with same resourceid")
	s.DeleteEntityByResourceID("k8snode", "cluster:ns:multinode")
	s.DeleteEntityByResourceID("k8snode", "cluster:ns:multinode2")
}

func TestCreateRelByUid(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	q := NewQueryService(dc)
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "namespace", Type: "uid"})
	dc.CreateSchema(db.Schema{Predicate: "cluster", Type: "uid"})

	podMeta := `{
		"name": "pod",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "labels",
				"fieldtype": "json",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "objtype",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "namespace",
				"fieldtype": "relationship",
				"refdatatype": "namespace",
				"mandatory": true,
				"index": false,
				"cardinality": "one"
			},
			{
				"fieldname": "cluster",
				"fieldtype": "relationship",
				"refdatatype": "cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "one"
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
		"name": "namespace",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			},
			{
				"fieldname": "cluster",
				"fieldtype": "relationship",
				"refdatatype": "cluster",
				"mandatory": true,
				"index": false,
				"cardinality": "one"
			}
		]
	}`
	nsMap := make(map[string]interface{})
	json.Unmarshal([]byte(nsMeta), &nsMap)
	s.CreateEntity("metadata", nsMap)
	clusterMeta := `{
		"name": "cluster",
        "objtype" : "metadata",
		"fields": [
			{
				"fieldname": "name",
				"fieldtype": "string",
				"mandatory": true,
				"index": true,
				"cardinality": "one"
			}
			
		]
	}`
	clMap := make(map[string]interface{})
	json.Unmarshal([]byte(clusterMeta), &clMap)
	s.CreateEntity("metadata", clMap)

	list := []string{"pod", "cluster", "namespace"}
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

	// create namespace
	ns := map[string]interface{}{"name": "ns01", "cluster": "c01", "objtype": "namespace"}
	uid, err := s.CreateEntity("namespace", ns)
	defer s.dbclient.DeleteEntity(uid)

	// create pod data
	pod := `{
		"name": "pod01",
        "objtype": "pod",
        "cluster": "c01",
        "namespace": {
            "uid":"` + uid +
		`"}
	}`
	podMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(pod), &podMap)
	if err != nil {
		panic(err)
	}
	uid2, err := s.CreateEntity("pod", podMap)
	if err != nil {
		panic(err)
	}

	pods, err := s.GetEntity(uid2)
	if err != nil {
		assert.Fail(t, "Failed to get created pod")
	}
	o2 := pods["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, uid, o2["namespace"].([]interface{})[0].(map[string]interface{})["uid"], "pod's namespace should be linked to existing one")
	s.dbclient.DeleteEntity(uid2)
	s.DeleteEntityByResourceID("cluster", "cluster:c01")
}

func TestEntityUpdate(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
	pod := map[string]interface{}{
		"name":       "pod03",
		"resourceid": "pod:pod03",
		"objtype":    "pod",
	}
	pid, _ := s.CreateEntity("pod", pod)
	s.UpdateEntity(pid, map[string]interface{}{"name": "pod04"})
	pods, _ := s.GetEntity(pid)
	o := pods["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "pod04", o["name"], "pod name should be updated")
	// try update with lower resource version
	s.UpdateEntity(pid, map[string]interface{}{"name": "pod05", "resourceversion": "0"})
	pods, _ = s.GetEntity(pid)
	o = pods["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "pod04", o["name"], "pod should not be updated due to version conflict")
	ns := map[string]interface{}{
		"name":       "ns01",
		"resourceid": "ns:ns01",
		"objtype":    "namespace",
	}
	nid, _ := s.CreateEntity("namespace", ns)
	// update with relationship
	s.UpdateEntity(pid, map[string]interface{}{"ns": map[string]string{"uid": nid}})
	pods, _ = s.GetEntity(pid)
	o = pods["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, nid, o["ns"].([]interface{})[0].(map[string]interface{})["uid"], "pod name should be updated with edge")
	ns2 := map[string]interface{}{
		"name":       "ns02",
		"resourceid": "ns:ns02",
		"objtype":    "namespace",
	}
	nid2, _ := s.CreateEntity("namespace", ns2)
	s.UpdateEntity(pid, map[string]interface{}{"ns": map[string]string{"uid": nid2}})
	pods, _ = s.GetEntity(pid)
	o = pods["objects"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, nid2, o["ns"].([]interface{})[0].(map[string]interface{})["uid"], "pod name should be updated with edge")
	defer s.DeleteEntity(pid)
	defer s.DeleteEntity(nid)
	defer s.DeleteEntity(nid2)
}
