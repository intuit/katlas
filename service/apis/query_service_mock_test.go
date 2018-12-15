package apis

import (
	"encoding/json"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/intuit/katlas/service/db"
	"github.com/stretchr/testify/assert"
)

var pod01 map[string]interface{}

func TestMockGetQueryResultByKeyValue(t *testing.T) {
	log.Infoln("Start TestMockGetQueryResultByKeyValue ...")
	theDBMock := db.MockDGClient{} //create the mock

	//Create query map
	qm := map[string][]string{
		"name":  {"pod01"},
		"_type": {"K8sPod"},
		"ip":    {"172.20.32.128"},
	}

	//Response data
	byt := []byte(
		`{
			"objects": [
			{
			  "uid": "1x0002",
			  "resourceVersion": "6365014",
			  "startTime": "2018-09-01T10:01:03Z",
			  "resourceId": "unique_id_of_pod",
			  "status": "Running",
			  "description": "describe this object",
			  "ip": "172.20.32.128",
			  "name": "pod01",
			  "_type": "K8sPod"
			}
		  ]
		}`)

	if err := json.Unmarshal(byt, &pod01); err != nil {
		panic(err)
	}

	theDBMock.On("GetQueryResult", qm).Return(pod01, nil)
	queryService := QueryService{&theDBMock}
	p, _ := queryService.GetQueryResult(qm)
	o1 := p["objects"].([]interface{})[0].(map[string]interface{})
	// assert
	assert.Equal(t, o1["name"], "pod01", "pod01 name not the same as input")
	assert.Equal(t, o1["_type"], "K8sPod", "pod01 type not the same as input")
	assert.Equal(t, o1["ip"], "172.20.32.128", "pod01 ip not the same as input")
	assert.Equal(t, o1["resourceId"], "unique_id_of_pod", "pod01 resourceId not the same as input")

	theDBMock.AssertNumberOfCalls(t, "GetQueryResult", 1)
	theDBMock.AssertExpectations(t)
}

func TestMockGetQueryResultByKeyword(t *testing.T) {
	log.Infoln("Start TestMockGetQueryResultByKeyword ...")
	theDBMock := db.MockDGClient{} //create the mock

	//Create query map
	qm := map[string][]string{
		"keyword": {"172.20.32.128"},
	}

	//Response data
	byt := []byte(
		`{
			"objects": [
			{
			  "uid": "1x0002",
			  "resourceVersion": "6365014",
			  "startTime": "2018-09-01T10:01:03Z",
			  "resourceId": "unique_id_of_pod",
			  "status": "Running",
			  "description": "describe this object",
			  "ip": "172.20.32.128",
			  "name": "pod01",
			  "_type": "K8sPod"
			}
		  ]
		}`)

	if err := json.Unmarshal(byt, &pod01); err != nil {
		panic(err)
	}

	theDBMock.On("GetQueryResult", qm).Return(pod01, nil)
	queryService := QueryService{&theDBMock}
	p, _ := queryService.GetQueryResult(qm)
	o1 := p["objects"].([]interface{})[0].(map[string]interface{})
	// assert
	assert.Equal(t, o1["name"], "pod01", "pod01 name not the same as input")
	assert.Equal(t, o1["_type"], "K8sPod", "pod01 type not the same as input")
	assert.Equal(t, o1["ip"], "172.20.32.128", "pod01 ip not the same as input")
	assert.Equal(t, o1["resourceId"], "unique_id_of_pod", "pod01 resourceId not the same as input")

	theDBMock.AssertNumberOfCalls(t, "GetQueryResult", 1)
	theDBMock.AssertExpectations(t)
}
