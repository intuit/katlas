package serviceapitests

import (
	"testing"
)

type Query interface {
	KeywordSearch(t *testing.T, url string) (statusCode int, respBody []byte, err error)
	QslSearch(t *testing.T, url string) (statusCode int, respBody []byte, err error)
}

// Keyword search for query
func KeywordSearch(t *testing.T, url string) (statusCode int, respBody []byte, err error) {

	respStatusCode, respBody, _ := GetResponseByGet(t, url)
	return respStatusCode, respBody, nil
}

// Qsl search for query
func QslSearch(t *testing.T, url string) (statusCode int, respBody []byte, err error) {

	respStatusCode, respBody, _ := GetResponseByGet(t, url)
	return respStatusCode, respBody, nil
}

// Test keyword search query
func TestKeywordSearch(t *testing.T) {
	testURL1 := TestBaseURL + "/v1/query?keyword=node01"
	testURL2 := TestBaseURL + "/v1/query?keyword=node02"

	tests := []TestStruct{
		{testURL1, "", 200, "", 0},
		{testURL2, "", 200, "", 0},
	}

	for i, testCase := range tests {
		resCode, resBody, _ := KeywordSearch(t, testCase.testURL)
		tests[i].observedStatusCode = resCode
		tests[i].responseBody = string(resBody)
	}
	DisplayTestCaseResults("TestKeywordSearch", tests, t, "uid")
}

// Test Qsl search query
func TestQslSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestQSLSearch in short mode")
	}
	testURL1 := TestBaseURL + "/v1/qsl/node[@name=\"node01\"]{*}"
	testURL2 := TestBaseURL + "/v1/qsl/node[@name=\"node01\"]{@labels}"

	tests := []TestStruct{
		{testURL1, "", 200, "", 0},
		{testURL2, "", 200, "", 0},
	}

	for i, testCase := range tests {
		resCode, resBody, _ := QslSearch(t, testCase.testURL)
		tests[i].observedStatusCode = resCode
		tests[i].responseBody = string(resBody)
	}
	DisplayTestCaseResults("TestQslSearch", tests, t, "uid")
}
