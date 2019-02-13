package serviceapitests

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

type TestStruct struct {
	testURL            string
	requestBody        string
	expectedStatusCode int
	responseBody       string
	observedStatusCode int
}

var TestBaseURL string

func TestMain(m *testing.M) {
	// parse and print command line flags
	flag.Parse()
	log.Printf("TestEnv=%s", TestConfig.TestEnv)
	log.Printf("Protocal=%s", TestConfig.Protocal)
	log.Printf("BasePath=%s", TestConfig.BasePath)
	log.Printf("Port=%s", TestConfig.Port)

	TestBaseURL = TestConfig.Protocal + "://" + TestConfig.BasePath + ":" + TestConfig.Port
	log.Printf("TestBaseUrl=%s", TestBaseURL)
	os.Exit(m.Run())
}

// GetResponseByGet ... Get response by GET method
func GetResponseByGet(t *testing.T, requestURL string) (statusCode int, respBody []byte, err error) {

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		t.Error(err)
	}

	respStatusCode, respBody, _ := getResponse(t, req)
	return respStatusCode, respBody, nil
}

// GetResponseByPost ... Get response by POST method
func GetResponseByPost(t *testing.T, requestURL string, requestBody string) (statusCode int, respBody []byte, err error) {

	jsonStr := []byte(requestBody)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error(err)
	}

	respStatusCode, respBody, _ := getResponse(t, req)
	return respStatusCode, respBody, nil
}

func getResponse(t *testing.T, req *http.Request) (statusCode int, responseBody []byte, err error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, body, nil
}

// DisplayTestCaseResults ... Compare testcase expected statuscode and observed statuscode to assert test success|fail
func DisplayTestCaseResults(functionalityName string, tests []TestStruct, t *testing.T) {

	for _, test := range tests {

		if test.observedStatusCode == test.expectedStatusCode {
			t.Logf("*** %s Passed Case: ***\n  -testURL : %s \n -requestBody : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n", functionalityName, test.testURL, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode)
		} else {
			t.Errorf("*** %s Failed Case: ***\n  -testURL : %s \n -requestBody : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n", functionalityName, test.testURL, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode)
		}
	}
}
