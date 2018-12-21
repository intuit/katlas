package main

import (
	//"encoding/json"

	"fmt"

	"github.com/gorilla/mux"

	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	//"github.com/dgraph-io/dgo/y"
)

var inCluster = false

type baseStruct struct {
	Name         string `json:"name"`
	Objtype      string `json:"objtype,omitempty"`
	Cluster      string `json:"cluster,omitempty"`
	Objnamespace string `json:"objnamespace,omitempty"`
}

type objectList struct {
	Me []map[string]string `json:"me"`
}

// Create function
func Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}

	var obj interface{}
	err2 := json.Unmarshal(body, &obj)
	if err2 != nil {
		log.Error(err2)
	}

	log.Infof("received %s: %+v", obj.(map[string]interface{})["objtype"].(string), obj)

	w.Write([]byte("Event " + obj.(map[string]interface{})["objtype"].(string) + " " + " received"))

}

// Health check
func Health(w http.ResponseWriter, r *http.Request) {
	log.Info("RestService is still running")
	w.Write([]byte("RestService is still running"))

}

// Sync function
func Sync(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}

	info := baseStruct{}
	err2 := json.Unmarshal(body, &info)
	if err2 != nil {
		log.Error(err2)
	}

	log.Info("sync received: ", info)
	w.Write([]byte("received sync for " + info.Objtype + "/" + info.Name))

}

// Restart function
func Restart(w http.ResponseWriter, r *http.Request) {
	// takes an object with Objtype restart and name as the cluster name
	// deletes the cluster from the database
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}

	info := baseStruct{}
	err2 := json.Unmarshal(body, &info)
	if err2 != nil {
		log.Error(err2)
	}

	w.Write([]byte("received restart for " + info.Objtype + "/" + info.Name))

}

// Up function
func Up(w http.ResponseWriter, r *http.Request) {
	log.Info("Up")
	w.Write([]byte("Up"))

}

// EntityHandler function
func EntityHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "error with data received", http.StatusInternalServerError)
			return
		}
	}()

	log.Infof("%+v\n", r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("error reading body")
		log.Error(err)
		return
	}

	var info map[string]interface{}
	err2 := json.Unmarshal(body, &info)
	if err2 != nil {
		log.Error("error with ", body, string(body))
		log.Error(err2)
		http.Error(w, "Failed to convert to JSON output", http.StatusInternalServerError)
		return
	}
	log.Infof("successfully received %+v", info)

	w.Write([]byte("object metadata for " + info["name"].(string) + " successfully received in entity"))

}

// QueryHandler ...
func QueryHandler(w http.ResponseWriter, r *http.Request) {

}

// SyncHandler ...
func SyncHandler(w http.ResponseWriter, r *http.Request) {

	resp := objectList{
		Me: []map[string]string{
			{
				"resourceversion": "1",
				"name":            "test-deployment",
				"namespace":       "namespace1",
			},
			{
				"resourceversion": "1",
				"name":            "test-ingress",
				"namespace":       "namespace1",
			},
			{
				"resourceversion": "1",
				"name":            "test-namespace",
				"namespace":       "namespace1",
			},
			{
				"resourceversion": "1",
				"name":            "test-pod",
				"namespace":       "namespace1",
			},

			{
				"resourceversion": "1",
				"name":            "test-replicaset",
				"namespace":       "namespace1",
			},
			{
				"resourceversion": "1",
				"name":            "test-service",
				"namespace":       "namespace1",
			},

			{
				"resourceversion": "1",
				"name":            "test-statefulset",
				"namespace":       "namespace1",
			},
		},
	}

	pdb, err3 := json.Marshal(resp)

	if err3 == nil {
		log.Info("sync sending response: ", pdb)
		w.Write(pdb)
	} else {
		fmt.Println(err3)
		log.Error(err3)
		w.Write([]byte("Error syncing: " + err3.Error()))
	}

}

// DeleteHandler ...
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "error with data received", http.StatusInternalServerError)
			return
		}
	}()
	vars := mux.Vars(r)
	meta := vars["metadata"]
	rid := vars["resourceid"]

	w.Write([]byte("object metadata for " + meta + "/" + rid + " successfully received for delete"))

}

func serve() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/entity/{metadata}", EntityHandler).Methods("GET", "POST", "DELETE")
	router.HandleFunc("/v1/query", QueryHandler).Methods("GET", "POST")
	router.HandleFunc("/v1/sync", SyncHandler).Methods("GET", "POST")
	router.HandleFunc("/v1/entity/{metadata}/{resourceid}", DeleteHandler).Methods("DELETE")
	router.HandleFunc("/health", Health).Methods("GET")
	log.Infof("Service started on port 8011")
	if inCluster {
		log.Error(http.ListenAndServeTLS(":8011", "server.crt", "server.key", router))
	} else {
		log.Error(http.ListenAndServe(":8011", router))
	}

}

func main() {
	serve()

}
