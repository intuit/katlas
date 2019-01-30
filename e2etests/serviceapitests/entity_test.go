package serviceapitests

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type Entity interface {
	CreateEntity(t *testing.T) (statusCode int, err error)
	//DeleteEntity(t *testing.T) (responseBody string, err error)
}

type TestStruct struct {
	requestBody        string
	expectedStatusCode int
	responseBody       string
	observedStatusCode int
}

//Suite tests all the functionality that Entity should implement
func CreateEntity(t *testing.T, url string, reqBody string) (statusCode int, respBody []byte, err error) {
	jsonStr := []byte(reqBody)
	//req, err := http.NewRequest("POST", "http://127.0.0.1:8011/v1/entity/K8sNode", bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	//log.Printf("Response code = %d", resp.StatusCode)
	//log.Printf("Response header content-type = %s", resp.Header.Get("Content-Type"))
	//log.Printf("Response body = %s", body)
	return resp.StatusCode, body, nil
}

//Test create entities
/*func TestCreateEntity(t *testing.T) {
	expected := 200
	res, _ := CreateEntity(t)
	if res != expected {
		t.Errorf("Got result: %d, want %d", res, expected)
		t.Fail()
	}
}*/

func TestCreateEntity(t *testing.T) {
	TestURL := TestBaseURL + "/v1/entity/K8sNode"
	log.Printf("TestURL= %s", TestURL)
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
		{node01, 200, "", 0},
		{node02, 200, "", 0},
	}

	for i, testCase := range tests {
		node := testCase.requestBody
		resCode, resBody, _ := CreateEntity(t, TestURL, node)
		/*if res != testCase.expectedStatusCode {
			t.Errorf("Got result: %d, want %d", res, tests[i].expectedStatusCode)
			t.Fail()
		}*/
		tests[i].observedStatusCode = resCode
		tests[i].responseBody = string(resBody)
	}
	DisplayTestCaseResults("TestCreateEntity", tests, t)
}

func DisplayTestCaseResults(functionalityName string, tests []TestStruct, t *testing.T) {

	for _, test := range tests {

		if test.observedStatusCode == test.expectedStatusCode {
			t.Logf("*** %s Passed Case: ***\n  -request body : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n", functionalityName, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode)
		} else {
			t.Errorf("*** %s Failed Case: ***\n  -request body : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n", functionalityName, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode)
			t.Fail()
		}
	}
}
