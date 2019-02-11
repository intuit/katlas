package serviceapitests

import (
	"testing"
)

type Entity interface {
	CreateEntity(t *testing.T, url string, reqBody string) (statusCode int, respBody []byte, err error)
}

// CreateEntity ... To create entity
func CreateEntity(t *testing.T, url string, reqBody string) (statusCode int, respBody []byte, err error) {

	respStatusCode, respBody, _ := GetResponseByPost(t, url, reqBody)
	return respStatusCode, respBody, nil
}

func TestCreateEntity(t *testing.T) {
	testURL1 := TestBaseURL + "/v1/entity/K8sNode"
	testURL2 := TestBaseURL + "/v1/entity/K8sNode"

	node01 := `{
		"objtype": "k8snode",
		"name": "node01",
		"labels": "testingnode01"
	  }`

	node02 := `{
	"objtype": "k8snode",
	"name": "node02",
	"labels": "testingnode02"
	}`

	tests := []TestStruct{
		{testURL1, node01, 200, "", 0},
		{testURL2, node02, 200, "", 0},
	}

	for i, testCase := range tests {
		resCode, resBody, _ := CreateEntity(t, testCase.testURL, testCase.requestBody)
		tests[i].observedStatusCode = resCode
		tests[i].responseBody = string(resBody)
	}
	DisplayTestCaseResults("TestCreateEntity", tests, t)
}
