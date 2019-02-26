package resources

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	metrics "github.com/intuit/katlas/service/metrics"
	"github.com/intuit/katlas/service/util"
	"io/ioutil"
	"net/http"
	"time"
)

// EntityGetHandlerV1_1 REST API for get Entity
func (s ServerResource) EntityGetHandlerV1_1(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	uid := vars[util.UID]

	start := time.Now()
	code := http.StatusOK
	defer func() {
		metrics.DgraphGetEntityLatencyHistogram.WithLabelValues(fmt.Sprintf("%d", code)).Observe(time.Since(start).Seconds())
	}()

	obj, err := s.EntitySvc.GetEntity(uid)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		code = http.StatusInternalServerError
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	// object not found
	if len(obj) == 0 {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr4xx.Inc()
		code = http.StatusNotFound
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"entity with id %s not found\"}", http.StatusNotFound, uid)))
		return
	}
	obj["status"] = http.StatusOK
	ret, _ := json.Marshal(obj)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

// MetaGetHandlerV1_1 REST API for get metadata
func (s ServerResource) MetaGetHandlerV1_1(w http.ResponseWriter, r *http.Request) {
	s.MetaGetHandler(w, r)
}

// MetaDeleteHandlerV1_1 REST API for delete metadata
func (s ServerResource) MetaDeleteHandlerV1_1(w http.ResponseWriter, r *http.Request) {
	s.MetaDeleteHandler(w, r)
}

// EntityDeleteHandlerV1_1 REST API for delete Entity
func (s ServerResource) EntityDeleteHandlerV1_1(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	uid := vars[util.UID]

	start := time.Now()
	code := http.StatusOK
	defer func() {
		metrics.DgraphDeleteEntityLatencyHistogram.WithLabelValues(fmt.Sprintf("%d", code)).Observe(time.Since(start).Seconds())
	}()

	err := s.EntitySvc.DeleteEntity(uid)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		code = http.StatusInternalServerError
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	msg := map[string]interface{}{
		"status": http.StatusOK,
		"objects": []map[string]interface{}{
			{
				"uid": uid,
			},
		},
	}
	ret, _ := json.Marshal(msg)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

// EntityCreateHandlerV1_1 REST API for create Entity
func (s ServerResource) EntityCreateHandlerV1_1(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	clusterName := r.Header.Get(util.ClusterName)
	metas, ok := r.URL.Query()[util.ObjType]
	if !ok || len(metas[0]) < 1 {
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"metadata not found from parameters\"}", http.StatusBadRequest)))
		return
	}
	meta := metas[0]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	payload, err := buildEntityData(clusterName, meta, body, false)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr4xx.Inc()
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
		return
	}

	start := time.Now()
	code := http.StatusOK
	defer func() {
		metrics.DgraphCreateEntityLatencyHistogram.WithLabelValues(fmt.Sprintf("%d", code)).Observe(time.Since(start).Seconds())
	}()

	uid, err := s.EntitySvc.CreateEntity(meta, payload.(map[string]interface{}))
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		code = http.StatusInternalServerError
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	msg := map[string]interface{}{
		"status": http.StatusOK,
		"objects": []map[string]interface{}{
			{
				"uid":     uid,
				"objtype": meta,
			},
		},
	}
	ret, _ := json.Marshal(msg)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

// EntityUpdateHandlerV1_1 REST API for update Entity
func (s ServerResource) EntityUpdateHandlerV1_1(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	uuid := vars[util.UID]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	payload := make(map[string]interface{}, 0)
	err = json.Unmarshal(body, &payload)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr4xx.Inc()
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
		return
	}

	start := time.Now()
	code := http.StatusOK
	defer func() {
		metrics.DgraphUpdateEntityLatencyHistogram.WithLabelValues(fmt.Sprintf("%d", code)).Observe(time.Since(start).Seconds())
	}()

	err = s.EntitySvc.UpdateEntity(uuid, payload)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		code = http.StatusInternalServerError
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	msg := map[string]interface{}{
		"status": http.StatusOK,
		"objects": []map[string]interface{}{
			{
				"uid": uuid,
			},
		},
	}
	ret, _ := json.Marshal(msg)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

// EntitySyncHandlerV1_1 REST API to sync entities
func (s ServerResource) EntitySyncHandlerV1_1(w http.ResponseWriter, r *http.Request) {
	s.EntitySyncHandler(w, r)
}

// QueryHandlerV1_1 REST API for get Query Response
func (s ServerResource) QueryHandlerV1_1(w http.ResponseWriter, r *http.Request) {
	s.QueryHandler(w, r)
}

// MetaCreateHandlerV1_1 REST API for create Metadata
func (s ServerResource) MetaCreateHandlerV1_1(w http.ResponseWriter, r *http.Request) {
	s.MetaCreateHandler(w, r)
}

// MetaUpdateHandlerV1_1 REST API for update Metadata
func (s ServerResource) MetaUpdateHandlerV1_1(w http.ResponseWriter, r *http.Request) {
	s.MetaUpdateHandler(w, r)
}

// SchemaUpsertHandlerV1_1 REST API for create Schema
func (s ServerResource) SchemaUpsertHandlerV1_1(w http.ResponseWriter, r *http.Request) {
	s.SchemaUpsertHandler(w, r)
}

// SchemaDropHandlerV1_1 remove db schema
func (s ServerResource) SchemaDropHandlerV1_1(w http.ResponseWriter, r *http.Request) {
	s.SchemaDropHandler(w, r)
}

// QSLHandlerV1_1 handles requests for QSL
func (s *ServerResource) QSLHandlerV1_1(w http.ResponseWriter, r *http.Request) {
	s.QSLHandler(w, r)
}
