package db

import (
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"
)

func TestDGClient(t *testing.T) {
	client := NewDGClient("127.0.0.1:9080")
	defer client.Close()
	// create node
	node := map[string]interface{}{
		"_type":  "K8sNode",
		"name":   "node01",
		"labels": "testingnode",
	}
	nids, _ := client.CreateEntity("K8sNode", node)
	var nid string
	for _, v := range nids {
		nid = v
		break
	}

	defer cleanUP(client, nids)

	// create pod
	pod := map[string]interface{}{
		"_type":           "K8sPod",
		"description":     "describe this object",
		"resourceid":      "unique_id_of_pod",
		"name":            "pod01",
		"resourceversion": "6365014",
		"startTime":       "2018-09-01T10:01:03Z",
		"status":          "Running",
		"ip":              "172.20.32.128",
	}
	uids, err := client.CreateEntity("K8sPod", pod)
	if err != nil {
		log.Fatalf("create testing data pod01 failed %v", err)
	}
	defer cleanUP(client, uids)
	for _, v := range uids {
		// get pod01 by uid
		pod01, err := client.GetEntity("K8sPod", v)
		if err != nil {
			log.Fatalf("failed to get pod01 entity %v", err)
		}
		o2 := pod01["objects"].([]interface{})[0].(map[string]interface{})
		// assert resourceid is the same as input
		assert.Equal(t, o2["resourceid"], "unique_id_of_pod", "pod01 resourceid not the same as input")
		// create relationship
		client.CreateOrDeleteEdge("K8sPod", v, "K8sNode", nid, "runsOn", create)
		// get pod again to check rel
		pod01, _ = client.GetEntity("K8sPod", v)
		o3 := pod01["objects"].([]interface{})[0].(map[string]interface{})
		if val, ok := o3["runsOn"]; ok {
			rel := val.([]interface{})[0].(map[string]interface{})
			// assert pod01 is runsOn node01
			assert.Equal(t, rel["name"], "node01", "pod01 doesn't runsOn expected node01")
		}
		// update pod01 status to Failed
		update := make(map[string]interface{})
		update["status"] = "Failed"
		client.UpdateEntity("K8sPod", v, update)
		pod01, _ = client.GetEntity("K8sPod", v)
		o4 := pod01["objects"].([]interface{})[0].(map[string]interface{})
		assert.Equal(t, o4["status"], "Failed", "pod01 status not update to Failed")
		// remove edges from pod01 to node01
		client.CreateOrDeleteEdge("K8sPod", v, "K8sNode", nid, "runsOn", delete)
		pod01, _ = client.GetEntity("K8sPod", v)
		o5 := pod01["objects"].([]interface{})[0].(map[string]interface{})
		val := o5["runsOn"]
		if val != nil {
			log.Fatalf("relationship still exist after call delete edge API")
		}
		client.CreateOrDeleteEdge("K8sPod", v, "K8sNode", nid, "runsOn", create)
		client.SetFieldToNull(map[string]interface{}{"runsOn": nil, "uid": v})
		pod01, _ = client.GetEntity("K8sPod", v)
		o6 := pod01["objects"].([]interface{})[0].(map[string]interface{})
		_, ok := o6["runsOn"]
		if ok {
			log.Fatalf("relationship still exist after call set field to null API")
		}

	}
}

func TestCreateIndex(t *testing.T) {
	client := NewDGClient("127.0.0.1:9080")
	s := Schema{Predicate: "testindex", Type: "string", Count: true, List: true, Index: true,
		Upsert: true, Tokenizer: []string{"hash", "fulltext"},
	}
	err := client.CreateSchema(s)
	assert.Nil(t, err)
	client.DropSchema("testindex")
}

func TestGetSchemaFromDB(t *testing.T) {
	client := NewDGClient("127.0.0.1:9080")
	defer client.Close()
	smds, err := client.GetSchemaFromDB()
	if err != nil {
		log.Fatalf("failed to get schema db %v", err)
	}
	log.Infof("From TestGetSchemaFromDB returned: [%v]\n", smds)
	assert.Nil(t, err)
}

func TestGetSchemaFromCache(t *testing.T) {
	client := NewDGClient("127.0.0.1:9080")
	defer client.Close()

	//Creates an LRU cache of the given size
	var err error
	LruCache, err = lru.New(5)
	if err != nil {
		log.Errorf("err: %v", err)
	}
	log.Infoln("LRU cache created with given size")
	InitLruCacheDBSchema = false
	//Get schema with fetching db
	smds, err := client.GetSchemaFromCache(LruCache)
	if err != nil {
		log.Fatalf("failed to get schema from LruCache %v", err)
	}
	if reflect.TypeOf(smds).Kind() == reflect.Slice {
		log.Debugln("The smds is a slice.")
	}
	log.Infof("From TestGetSchemaFromCache returned from LruCache with fetching db: [%v]\n", smds)
	assert.Nil(t, err)

	//Get schema without fetching db
	smds, err = client.GetSchemaFromCache(LruCache)
	if err != nil {
		log.Fatalf("failed to get schema from LruCache %v", err)
	}
	log.Infof("From TestGetSchemaFromCache returned from LruCache without fetching db: [%v]\n", smds)
	assert.Nil(t, err)

	//Remove schema from cache
	client.RemoveDBSchemaFromCache(LruCache)

	//Get schema again should fetch DB again
	smds, err = client.GetSchemaFromCache(LruCache)
	if err != nil {
		log.Fatalf("failed to get schema from LruCache %v", err)
	}

	log.Infof("From TestGetSchemaFromCache returned from LruCache second time: [%v]\n", smds)
	assert.Nil(t, err)
}

func cleanUP(client *DGClient, uids map[string]string) {
	for _, v := range uids {
		err := client.DeleteEntity(v)
		if err != nil {
			log.Printf("something wrong to delete entity %s " + v)
		}
	}
}
