package apis

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"
	lru "github.com/hashicorp/golang-lru"
	"github.com/intuit/katlas/service/db"
	"github.com/stretchr/testify/assert"
)

func TestGetQueryResultByKeyValue(t *testing.T) {

	//Expected values
	expectPodName := "pod01"
	expectObjectsKey := "objects"

	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			log.Errorf("Exception: [%v]\n", err)
			assert.Fail(t, "Exception: [%v]\n", err)
		}
	}()

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	ms := NewMetaService(dc)
	//create entity for query later
	puid := createPod(dc, ms)
	defer deletePod(dc, ms, puid)

	s := NewQueryService(dc)

	//Create query map
	m := map[string][]string{
		"name":    {"pod01"},
		"objtype": {"Pod"},
		"ip":      {"172.20.32.128"},
	}

	//Get QueryResult
	qr := make(map[string]interface{})
	qr, err := s.GetQueryResult(m)
	if err != nil {
		log.Errorf("Query result err: [%v]\n", err)

	}
	log.Infof("Query result: [%v]\n", qr)

	if value, ok := qr[expectObjectsKey]; ok {
		log.Debugln("value: ", value)
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			log.Debugln("The interface is a slice.")
		}
		if value != nil {
			o := qr[expectObjectsKey].([]interface{})[0].(map[string]interface{})
			if val, ok := o["name"]; ok {
				assert.Equal(t, expectPodName, val, "Pod name is not equal to expected.")
			} else {
				assert.Fail(t, "Failed to GetQueryResult")
			}
		} else {
			assert.Fail(t, "Returned query result is empty!")
		}
	} else {
		assert.Fail(t, "map key is not found!")
	}
}

func TestGetQueryResultByKeywordSearch(t *testing.T) {
	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()

	var err error
	db.LruCache, err = lru.New(5)
	if err != nil {
		log.Errorf("err: %v", err)
	}
	log.Infoln("LRU cache created with given size")

	s := NewQueryService(dc)

	//Create query map
	m := map[string][]string{
		"keyword": {"pod"},
	}

	//Get QueryResult
	qr := make(map[string]interface{})
	qr, err = s.GetQueryResult(m)
	if err != nil {
		log.Errorf("Test Get Query result err: [%v]\n", err)

	}
	log.Infof("Test Get Query Result: [%v]\n", qr)
	assert.Nil(t, err)
}

func createPod(dc *db.DGClient, ms *MetaService) (uid string) {
	s := NewEntityService(dc)
	// create pod
	pod := map[string]interface{}{
		"objtype":         "Pod",
		"resourceId":      "unique_id_of_pod01",
		"name":            "pod01",
		"resourceversion": "6365014",
		"starttime":       "2018-09-01T10:01:03Z",
		"status":          "Running",
		"ip":              "172.20.32.128",
	}
	// create index for query
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term", "trigram"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term", "trigram"}})
	dc.CreateSchema(db.Schema{Predicate: "ip", Type: "string", Index: true, Tokenizer: []string{"term"}})

	pids, _ := s.CreateEntity("K8sPod", pod)
	var pid string
	for _, v := range pids {
		pid = v
		break
	}

	p, _ := s.GetEntity("K8sPod", pid)
	o := p["objects"].([]interface{})[0].(map[string]interface{})
	if val, ok := o["name"]; ok {
		log.Infof("name: [%v]\n", val)
	} else {
		log.Errorln("failed to create k8s pod object")
	}
	return pid
}

func deletePod(dc *db.DGClient, ms *MetaService, uid string) {
	s := NewEntityService(dc)
	s.DeleteEntity(uid)
}
