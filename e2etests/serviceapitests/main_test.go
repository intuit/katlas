package serviceapitests

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type TestStruct struct {
	testCaseName       string
	testURL            string
	requestBody        string
	expectedStatusCode int
	responseBody       string
	observedStatusCode int
}

var TestBaseURL string

func TestMain(m *testing.M) {
	// Parse and print command line flags
	flag.Parse()
	log.Printf("TestEnv=%s", TestConfig.TestEnv)
	log.Printf("Protocal=%s", TestConfig.Protocal)
	log.Printf("BasePath=%s", TestConfig.BasePath)
	log.Printf("Port=%s", TestConfig.Port)

	TestBaseURL = TestConfig.Protocal + "://" + TestConfig.BasePath + ":" + TestConfig.Port
	log.Printf("TestBaseUrl=%s", TestBaseURL)
	setupPreTestData()
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

func setupPreTestData() {
	testURL1 := TestBaseURL + "/v1/entity/node"
	testURL2 := TestBaseURL + "/v1/entity/node"

	node01 := GetStrFromJSONFile("entity_node01.json")
	log.Printf("node01 = %s", node01)

	node02 := GetStrFromJSONFile("entity_node02.json")
	log.Printf("node02 = %s", node02)

	tests := []TestStruct{
		{"CreateNode01", testURL1, node01, 200, "", 0},
		{"CreateNode02", testURL2, node02, 200, "", 0},
	}

	for i, testCase := range tests {
		jsonStr := []byte(testCase.requestBody)
		req, err := http.NewRequest("POST", tests[i].testURL, bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Response status code: %d \n Response body: %s \n", resp.StatusCode, body)
	}
}

//GetStrFromJsonFile ... Get test input data from a json file under test-fixtures directory
func GetStrFromJSONFile(fileStr string) (jsonStr string) {
	filepath := filepath.Join("test-fixtures", fileStr)
	log.Printf("filepath= %s", filepath)
	jsonFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Panicln(err)
	}
	return string(jsonFile)
}

// DisplayTestCaseResults ... Compare testcase expected statuscode with observed statuscode, and expected response string with observer response string to assert test success|fail
func DisplayTestCaseResults(functionalityName string, tests []TestStruct, t *testing.T, expectResponseStr string) {

	for _, test := range tests {

		if test.observedStatusCode == test.expectedStatusCode {
			if strings.Contains(expectResponseStr, "&") {
				expectResStrs := strings.Split(expectResponseStr, "&")
				if strings.Contains(test.responseBody, expectResStrs[0]) && strings.Contains(test.responseBody, expectResStrs[1]) {
					t.Logf("*** %s Passed Case: ***\n  -testCaseName : %s \n -testURL : %s \n -requestBody : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n -expectResponseStr : %s \n",
						functionalityName, test.testCaseName, test.testURL, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode, expectResponseStr)
				} else {
					t.Errorf("*** %s Failed Case: ***\n  -testCaseName : %s \n -testURL : %s \n -requestBody : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n -expectResponseStr : %s \n",
						functionalityName, test.testCaseName, test.testURL, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode, expectResponseStr)
				}
			} else if strings.Contains(expectResponseStr, "|") {
				expectResStrs := strings.Split(expectResponseStr, "|")
				if strings.Contains(test.responseBody, expectResStrs[0]) || strings.Contains(test.responseBody, expectResStrs[1]) {
					t.Logf("*** %s Passed Case: ***\n  -testCaseName : %s \n -testURL : %s \n -requestBody : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n -expectResponseStr : %s \n",
						functionalityName, test.testCaseName, test.testURL, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode, expectResponseStr)
				} else {
					t.Errorf("*** %s Failed Case: ***\n  -testCaseName : %s \n -testURL : %s \n -requestBody : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n -expectResponseStr : %s \n",
						functionalityName, test.testCaseName, test.testURL, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode, expectResponseStr)
				}
			} else if strings.Contains(test.responseBody, expectResponseStr) {
				t.Logf("*** %s Passed Case: ***\n  -testCaseName : %s \n -testURL : %s \n -requestBody : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n -expectResponseStr : %s \n",
					functionalityName, test.testCaseName, test.testURL, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode, expectResponseStr)
			} else {
				t.Errorf("*** %s Failed Case: ***\n  -testCaseName : %s \n -testURL : %s \n -requestBody : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n -expectResponseStr : %s \n",
					functionalityName, test.testCaseName, test.testURL, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode, expectResponseStr)
			}
		} else {
			t.Errorf("*** %s Failed Case: ***\n  -testCaseName : %s \n -testURL : %s \n -requestBody : %s \n -expectedStatus : %d \n -responseBody : %s \n -observedStatusCode : %d \n -expectResponseStr : %s \n",
				functionalityName, test.testCaseName, test.testURL, test.requestBody, test.expectedStatusCode, test.responseBody, test.observedStatusCode, expectResponseStr)
		}
	}
}
