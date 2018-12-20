package resources

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/intuit/katlas/service/apis"
	"github.com/intuit/katlas/service/util"
)

// ServerResource handle http request
type ServerResource struct {
	EntitySvc *apis.EntityService
	QuerySvc  *apis.QueryService
	MetaSvc   *apis.MetaService
	// TODO:
	// add metadata service, audit service and spec service after API ready
}

// EntityGetHandler REST API for get Entity
func (s ServerResource) EntityGetHandler(w http.ResponseWriter, r *http.Request) {
	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	meta := vars[util.Metadata]
	uid := vars[util.UID]
	obj, err := s.EntitySvc.GetEntity(meta, uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(ret)
}

// MetaGetHandler REST API for get Entity
func (s ServerResource) MetaGetHandler(w http.ResponseWriter, r *http.Request) {
	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	name := vars[util.Name]
	obj, err := s.MetaSvc.GetMetadata(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(ret)
}

// EntityDeleteHandler REST API for delete Entity
func (s ServerResource) EntityDeleteHandler(w http.ResponseWriter, r *http.Request) {
	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	meta := vars[util.Metadata]
	rid := vars[util.ResourceID]
	err := s.EntitySvc.DeleteEntityByResourceID(meta, rid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"deleted\": \"%s\", \"type\": \"%s\"}", rid, meta)))
}

// EntityCreateHandler REST API for create Entity
func (s ServerResource) EntityCreateHandler(w http.ResponseWriter, r *http.Request) {
	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	meta := vars[util.Metadata]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	var payload map[string]interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uids, err := s.EntitySvc.CreateEntity(meta, payload)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ret, err := json.Marshal(uids)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

// EntitySyncHandler REST API to sync entities
func (s ServerResource) EntitySyncHandler(w http.ResponseWriter, r *http.Request) {
	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	meta := vars[util.Metadata]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	var payload []map[string]interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.EntitySvc.SyncEntities(meta, payload)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"synced\": \"done\", \"type\": \"%s\"}", meta)))
}

// QueryHandler REST API for get Query Response
func (s ServerResource) QueryHandler(w http.ResponseWriter, r *http.Request) {
	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryMap := r.URL.Query()

	obj, err := s.QuerySvc.GetQueryResult(queryMap)
	if err != nil {
		http.Error(w, "Service Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "Failed to convert to JSON output", http.StatusInternalServerError)
		return
	}
	w.Write(ret)
}

// MetaCreateHandler REST API for create Metadata
func (s ServerResource) MetaCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	meta := vars[util.Metadata]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	var payload map[string]interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uids, err := s.MetaSvc.CreateMetadata(meta, payload)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ret, err := json.Marshal(uids)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(ret)
}

// TODO:
// Add more supported rest APIs
