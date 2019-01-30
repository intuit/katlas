package serviceapitests

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type Query interface {
	KeywordSearch(t *testing.T) (statusCode int, err error)
	NameValueSearch(t *testing.T) (responseBody string, err error)
}

//Suite tests all the functionality that Query should implement
func KeywordSearch(t *testing.T) (statusCode int, err error) {
	TestURL := TestBaseURL + "/v1/query?name=node01"
	log.Printf("TestURL= %s", TestURL)

	req, err := http.NewRequest("GET", TestURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	t.Logf("TestURL= %s", TestURL)
	log.Printf("Response code = %d", resp.StatusCode)
	log.Printf("Response header content-type = %s", resp.Header.Get("Content-Type"))
	log.Printf("Response body = %s", body)
	return resp.StatusCode, nil
}

func NameValueSearch(t *testing.T) (responseBody string, err error) {
	TestURL := TestBaseURL + "/v1/query?name=node01"
	log.Printf("TestURL= %s", TestURL)
	req, err := http.NewRequest("GET", TestURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	t.Logf("TestURL= %s", TestURL)
	log.Printf("Response code = %d", resp.StatusCode)
	log.Printf("Response header content-type = %s", resp.Header.Get("Content-Type"))
	log.Printf("Response body = %s", body)
	return string(body), nil
}

//Test one for KeywordSearch
func TestKeywordSearch(t *testing.T) {
	log.Println("******Testing TestKeywordSearch start ******")
	expected := 200
	res, _ := KeywordSearch(t)
	if res != expected {
		t.Errorf("Got result: %d, want %d", res, expected)
		t.Fail()
	}
}

func TestNameValueSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestNameValueSearch in short mode")
	}
	//expected := 200
	res, _ := NameValueSearch(t)
	t.Logf("Response body: %s", res)
	if res == "" {
		t.Error("Got result is empty")
		t.Fail()
	}
}
