package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	time "time"

	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//"reflect"
	//"github.com/dgraph-io/dgo/y"
	"net/http/httputil"
)

// AppNamespace used for only monitor assets in a specific namespace
// otherwise AppNamespace will refer to all namespaces that the controller has access to
var AppNamespace string

// ClusterName for running in a cluster
var ClusterName = os.Getenv("CLUSTER_NAME")

// RestSvcEndpoint the URL that k8s data will send to
var RestSvcEndpoint = os.Getenv("TARGET_URL")

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{}) error
	ObjectDeleted(obj interface{}, key string) error
	ObjectUpdated(objOld, objNew interface{}) error
}

// SendJSONQuery send requests to REST api
func SendJSONQuery(obj interface{}, url string) (int, []byte) {
	//url := "http://localhost:8011/create"

	s, err := json.Marshal(obj)
	if err != nil {
		log.Error("failed to marshal object in SendJSONQuery")
		log.Error(err)
	}
	log.Infof("sent to %s:\n    %s", url, string(s))
	payload := bytes.NewBuffer(s)

	req, _ := http.NewRequest("POST", url, payload)
	req.Close = true
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", os.Getenv("AUTH_HEADER"))
	req.Header.Add("clustername", ClusterName)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		return 404, nil
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	log.Infof("%s response %+v\n", url, res)
	log.Infof("%s body: %s\n", url, string(body))
	log.Info("*************************************************************************")
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Info(string(requestDump))
	return res.StatusCode, body
}

// SendDeleteRequest send request to delete k8s objects
func SendDeleteRequest(url string) (int, []byte) {
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Close = true
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", os.Getenv("AUTH_HEADER"))
	// Fetch Request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		return 404, nil
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof("%s response %+v\n", url, res)
	log.Infof("%s body: %s\n", url, string(body))
	log.Info("*************************************************************************")
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Info(string(requestDump))
	return res.StatusCode, body
}

// SendJSONQueryWithRetries retry requests if error occurs
func SendJSONQueryWithRetries(obj interface{}, url string) ([]byte, error) {
	// try sending the query 5 times and if it fails
	status, body := SendJSONQuery(obj, url)
	maxTries := 5
	cur := 0
	for status != 200 && cur < maxTries {
		time.Sleep(2000 * time.Millisecond)
		status, body = SendJSONQuery(obj, url)
		cur = cur + 1
	}
	fmt.Println()
	time.Sleep(100 * time.Millisecond)

	if status == 200 {
		return body, nil
	}
	return nil, errors.New("sending object failed too many times")

}

// GetKubernetesClient retrieve the Kubernetes cluster client from outside of the cluster
func GetKubernetesClient() kubernetes.Interface {
	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"
	var config *rest.Config
	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Infof("getClusterConfig: %v", err)

		// if the container is running inside the cluster, use the incluster config
		var err2 error
		config, err2 = rest.InClusterConfig()
		if err2 != nil {
			panic(err.Error())
		}
	}

	// generate the client based off of the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client
}
