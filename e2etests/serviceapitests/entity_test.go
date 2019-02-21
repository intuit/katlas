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
	testURL1 := TestBaseURL + "/v1/entity/node"
	testURL2 := TestBaseURL + "/v1/entity/node"

	node11 := `{
		"objtype": "node",
		"name": "node11",
		"labels": "testingnode11"
	  }`

	node12 := `{
	"objtype": "node",
	"name": "node12",
	"labels": "testingnode12"
	}`

	tests := []TestStruct{
		{"TestCreateNode11", testURL1, node11, 200, "", 0},
		{"TestCreateNode12", testURL2, node12, 200, "", 0},
	}

	for i, testCase := range tests {
		t.Run(testCase.testCaseName, func(t *testing.T) {
			resCode, resBody, _ := CreateEntity(t, testCase.testURL, testCase.requestBody)
			tests[i].observedStatusCode = resCode
			tests[i].responseBody = string(resBody)
		})
	}
	DisplayTestCaseResults("TestCreateEntity", tests, t, "blank-0|node")
}
