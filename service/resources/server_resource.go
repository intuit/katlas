package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/intuit/katlas/service/apis"
	"github.com/intuit/katlas/service/db"
	metrics "github.com/intuit/katlas/service/metrics"
	"github.com/intuit/katlas/service/util"
	"github.com/mitchellh/mapstructure"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/apps/v1beta2"
	core_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	"reflect"
	"strings"
)

// ServerResource handle http request
type ServerResource struct {
	EntitySvc *apis.EntityService
	QuerySvc  *apis.QueryService
	MetaSvc   *apis.MetaService
	QSLSvc    *apis.QSLService
	// TODO:
	// add metadata service, audit service and spec service after API ready
}

// EntityGetHandler REST API for get Entity
func (s ServerResource) EntityGetHandler(w http.ResponseWriter, r *http.Request) {

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

// MetaGetHandler REST API for get metadata
func (s ServerResource) MetaGetHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	name := strings.ToLower(vars[util.Name])
	obj, err := s.MetaSvc.GetMetadata(name)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	if obj != nil {
		metrics.KatlasNumReq2xx.Inc()
		ret := []byte(fmt.Sprintf("{\"status\": %v, \"objects\": [", http.StatusOK))
		meta, _ := json.Marshal(obj)
		ret = append(ret, meta...)
		ret = append(ret, []byte("]}")...)
		w.Write(ret)
		return
	}
	metrics.KatlasNumReqErr.Inc()
	metrics.KatlasNumReqErr4xx.Inc()
	w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"metadata %s not found\"}", http.StatusNotFound, name)))
}

// MetaDeleteHandler REST API for delete metadata
func (s ServerResource) MetaDeleteHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	name := vars[util.Name]
	err := s.MetaSvc.DeleteMetadata(name)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr4xx.Inc()
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusConflict, trim(err.Error()))))
		return
	}
	msg := map[string]interface{}{
		"status": http.StatusOK,
		"objects": []map[string]interface{}{
			{
				"name":    name,
				"objtype": "metadata",
			},
		},
	}
	ret, _ := json.Marshal(msg)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

// EntityDeleteHandler REST API for delete Entity
func (s ServerResource) EntityDeleteHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	meta := vars[util.Metadata]
	rid := vars[util.ResourceID]

	start := time.Now()
	code := http.StatusOK
	defer func() {
		metrics.DgraphDeleteEntityLatencyHistogram.WithLabelValues(fmt.Sprintf("%d", code)).Observe(time.Since(start).Seconds())
	}()

	err := s.EntitySvc.DeleteEntityByResourceID(meta, rid)
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
				"resourceid": rid,
				"objtype":    meta,
			},
		},
	}
	ret, _ := json.Marshal(msg)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

// EntityCreateHandler REST API for create Entity
func (s ServerResource) EntityCreateHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	meta := vars[util.Metadata]
	clusterName := r.Header.Get(util.ClusterName)

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

// EntityUpdateHandler REST API for update Entity
func (s ServerResource) EntityUpdateHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	meta := vars[util.Metadata]
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
				"uid":     uuid,
				"objtype": meta,
			},
		},
	}
	ret, _ := json.Marshal(msg)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

// EntitySyncHandler REST API to sync entities
func (s ServerResource) EntitySyncHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	meta := vars[util.Metadata]
	clusterName := r.Header.Get(util.ClusterName)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	payload, err := buildEntityData(clusterName, meta, body, true)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr4xx.Inc()
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
		return
	}
	err = s.EntitySvc.SyncEntities(meta, payload.([]map[string]interface{}))
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	w.Write([]byte(fmt.Sprintf("{\"synced\": \"done\", \"type\": \"%s\"}", meta)))

	metrics.KatlasNumReq2xx.Inc()
}

// QueryHandler REST API for get Query Response
func (s ServerResource) QueryHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	queryMap := r.URL.Query()

	code := http.StatusOK
	start := time.Now()
	defer func() {
		metrics.KatlasQueryLatencyHistogram.WithLabelValues("katlas", "*", "None", "dev", "containers", "GET", fmt.Sprintf("%d", code), "/**").Observe(time.Since(start).Seconds())
	}()

	obj, err := s.QuerySvc.GetQueryResult(queryMap)
	if err != nil {
		code = http.StatusInternalServerError
		metrics.KatlasNumReqErr5xx.Inc()
		metrics.KatlasNumReqErr.Inc()
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	obj["status"] = http.StatusOK
	ret, _ := json.Marshal(obj)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

// MetaCreateHandler REST API for create Metadata
func (s ServerResource) MetaCreateHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	var payload interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr4xx.Inc()
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
		return
	}
	var msg map[string]interface{}
	if reflect.TypeOf(payload).Kind() == reflect.Slice {
		var rets []map[string]interface{}
		for _, p := range payload.([]interface{}) {
			uid, err := s.MetaSvc.CreateMetadata(p.(map[string]interface{}))
			if err != nil {
				metrics.KatlasNumReqErr.Inc()
				metrics.KatlasNumReqErr5xx.Inc()
				log.Error(err)
				w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
				return
			}

			rets = append(rets, map[string]interface{}{
				"uid":     uid,
				"objtype": p.(map[string]interface{})[util.Name],
			})
		}
		msg = map[string]interface{}{
			"status":  http.StatusOK,
			"objects": rets,
		}
		metrics.KatlasNumReq2xx.Inc()
	} else {
		uid, err := s.MetaSvc.CreateMetadata(payload.(map[string]interface{}))
		if err != nil {
			metrics.KatlasNumReqErr.Inc()
			metrics.KatlasNumReqErr5xx.Inc()
			log.Error(err)
			w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
			return
		}
		msg = map[string]interface{}{
			"status": http.StatusOK,
			"objects": []map[string]interface{}{
				{
					"uid":     uid,
					"objtype": payload.(map[string]interface{})[util.Name],
				},
			},
		}
		metrics.KatlasNumReq2xx.Inc()
	}
	ret, _ := json.Marshal(msg)
	w.Write(ret)
}

// MetaUpdateHandler REST API for update Metadata
func (s ServerResource) MetaUpdateHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	name := vars[util.Name]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	var payload interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr4xx.Inc()
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
		return
	}
	err = s.MetaSvc.UpdateMetadata(name, payload.(map[string]interface{}))
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	msg := map[string]interface{}{
		"status": http.StatusOK,
		"objects": []map[string]interface{}{
			{
				"name":    name,
				"objtype": "metadata",
			},
		},
	}

	ret, _ := json.Marshal(msg)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

// SchemaUpsertHandler REST API for create Schema
func (s ServerResource) SchemaUpsertHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	defer s.MetaSvc.RemoveSchemaCache(db.LruCache)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	var payload interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr4xx.Inc()
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
		return
	}
	var msg map[string]interface{}
	if reflect.TypeOf(payload).Kind() == reflect.Slice {
		var predicates []db.Schema
		err := mapstructure.Decode(payload, &predicates)
		if err != nil {
			metrics.KatlasNumReqErr.Inc()
			metrics.KatlasNumReqErr4xx.Inc()
			log.Error(err)
			w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
			return
		}
		names := make([]string, 0)
		for _, p := range predicates {
			err := s.MetaSvc.CreateSchema(p)
			if err != nil {
				log.Error(err)
				metrics.KatlasNumReqErr.Inc()
				metrics.KatlasNumReqErr5xx.Inc()
				w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
				return
			}
			names = append(names, p.Predicate)
		}
		msg = map[string]interface{}{
			"status":  http.StatusOK,
			"message": fmt.Sprintf("%v upsert successfully", names),
		}
		metrics.KatlasNumReq2xx.Inc()
	} else {
		var predicate db.Schema
		err := mapstructure.Decode(payload, &predicate)
		if err != nil {
			metrics.KatlasNumReqErr.Inc()
			metrics.KatlasNumReqErr4xx.Inc()
			log.Error(err)
			w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
			return
		}
		err = s.MetaSvc.CreateSchema(predicate)
		if err != nil {
			metrics.KatlasNumReqErr.Inc()
			metrics.KatlasNumReqErr5xx.Inc()
			log.Error(err)
			w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
			return
		}
		msg = map[string]interface{}{
			"status":  http.StatusOK,
			"message": fmt.Sprintf("%s upsert successfully", predicate.Predicate),
		}
		metrics.KatlasNumReq2xx.Inc()
	}
	ret, _ := json.Marshal(msg)
	w.Write(ret)
}

// SchemaDropHandler remove db schema
func (s ServerResource) SchemaDropHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	defer s.MetaSvc.RemoveSchemaCache(db.LruCache)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	predicate := vars[util.Name]

	err := s.MetaSvc.DropSchema(predicate)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	msg := map[string]interface{}{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("schema %s drop successfully", predicate),
	}
	ret, _ := json.Marshal(msg)
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

func buildEntityData(clusterName string, meta string, body []byte, isArray bool) (interface{}, error) {
	switch meta {
	case util.Namespace:
		if isArray {
			list := make([]map[string]interface{}, 0)
			data := []core_v1.Namespace{}
			err := json.Unmarshal(body, &data)
			if err != nil {
				return nil, err
			}
			for _, d := range data {
				namespace := map[string]interface{}{
					util.ObjType:         util.Namespace,
					util.Name:            d.ObjectMeta.Name,
					util.CreationTime:    d.ObjectMeta.CreationTimestamp,
					util.Cluster:         clusterName,
					util.ResourceVersion: d.ResourceVersion,
					util.K8sObj:          util.K8sObj,
					util.Labels:          d.ObjectMeta.GetLabels(),
					util.Asset:           getValues(&data, util.AssetID, "GetAnnotations"),
				}
				list = append(list, namespace)
			}
			return list, nil
		}
		data := core_v1.Namespace{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			util.ObjType:         util.Namespace,
			util.Name:            data.ObjectMeta.Name,
			util.CreationTime:    data.ObjectMeta.CreationTimestamp,
			util.Cluster:         clusterName,
			util.ResourceVersion: data.ResourceVersion,
			util.K8sObj:          util.K8sObj,
			util.Labels:          data.ObjectMeta.GetLabels(),
			util.Asset:           getValues(&data, util.AssetID, "GetAnnotations"),
		}, nil
	case util.Deployment:
		if isArray {
			list := make([]map[string]interface{}, 0)
			data := []v1beta2.Deployment{}
			err := json.Unmarshal(body, &data)
			if err != nil {
				return nil, err
			}
			for _, d := range data {
				deployment := map[string]interface{}{
					util.ObjType:           util.Deployment,
					util.Cluster:           clusterName,
					util.Name:              d.ObjectMeta.Name,
					util.CreationTime:      d.ObjectMeta.CreationTimestamp,
					util.Namespace:         d.ObjectMeta.Namespace,
					util.NumReplicas:       d.Spec.Replicas,
					util.AvailableReplicas: d.Status.AvailableReplicas,
					util.Strategy:          d.Spec.Strategy.Type,
					util.ResourceVersion:   d.ResourceVersion,
					util.Labels:            d.ObjectMeta.GetLabels(),
					util.K8sObj:            util.K8sObj,
				}
				// creata application from labels
				appList := createAppNameList(&d)
				if len(appList) > 0 {
					deployment[util.Application] = appList
				}
				list = append(list, deployment)
			}
			return list, nil
		}
		data := v1beta2.Deployment{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		deployment := map[string]interface{}{
			util.ObjType:           util.Deployment,
			util.Cluster:           clusterName,
			util.Name:              data.ObjectMeta.Name,
			util.CreationTime:      data.ObjectMeta.CreationTimestamp,
			util.Namespace:         data.ObjectMeta.Namespace,
			util.NumReplicas:       data.Spec.Replicas,
			util.AvailableReplicas: data.Status.AvailableReplicas,
			util.Strategy:          data.Spec.Strategy.Type,
			util.ResourceVersion:   data.ResourceVersion,
			util.Labels:            data.ObjectMeta.GetLabels(),
			util.K8sObj:            util.K8sObj,
		}
		// creata application from labels
		appList := createAppNameList(&data)
		if len(appList) > 0 {
			deployment[util.Application] = appList
		}
		return deployment, nil
	case util.Ingress:
		if isArray {
			list := make([]map[string]interface{}, 0)
			data := []ext_v1beta1.Ingress{}
			err := json.Unmarshal(body, &data)
			if err != nil {
				return nil, err
			}
			for _, d := range data {
				ingress := map[string]interface{}{
					util.ObjType:         util.Ingress,
					util.Cluster:         clusterName,
					util.Name:            d.ObjectMeta.Name,
					util.Namespace:       d.ObjectMeta.Namespace,
					util.CreationTime:    d.ObjectMeta.CreationTimestamp,
					util.DefaultBackend:  d.Spec.Backend,
					util.TSL:             d.Spec.TLS,
					util.Rules:           d.Spec.Rules,
					util.ResourceVersion: d.ObjectMeta.ResourceVersion,
					util.Labels:          d.ObjectMeta.GetLabels(),
					util.K8sObj:          util.K8sObj,
				}
				// creata application from labels
				appList := createAppNameList(&d)
				if len(appList) > 0 {
					ingress[util.Application] = appList
				}
				list = append(list, ingress)
			}
			return list, nil
		}
		data := ext_v1beta1.Ingress{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		ingress := map[string]interface{}{
			util.ObjType:         util.Ingress,
			util.Cluster:         clusterName,
			util.Name:            data.ObjectMeta.Name,
			util.Namespace:       data.ObjectMeta.Namespace,
			util.CreationTime:    data.ObjectMeta.CreationTimestamp,
			util.DefaultBackend:  data.Spec.Backend,
			util.TSL:             data.Spec.TLS,
			util.Rules:           data.Spec.Rules,
			util.ResourceVersion: data.ObjectMeta.ResourceVersion,
			util.Labels:          data.ObjectMeta.GetLabels(),
			util.K8sObj:          util.K8sObj,
		}
		appList := createAppNameList(&data)
		if len(appList) > 0 {
			ingress[util.Application] = appList
		}
		return ingress, nil
	case util.Pod:
		if isArray {
			list := make([]map[string]interface{}, 0)
			data := []core_v1.Pod{}
			err := json.Unmarshal(body, &data)
			if err != nil {
				return nil, err
			}
			for _, d := range data {
				pod := map[string]interface{}{
					util.ObjType:         util.Pod,
					util.Name:            d.ObjectMeta.Name,
					util.Namespace:       d.ObjectMeta.Namespace,
					util.CreationTime:    d.ObjectMeta.CreationTimestamp,
					util.Phase:           d.Status.Phase,
					util.NodeName:        d.Spec.NodeName,
					util.IP:              d.Status.PodIP,
					util.Containers:      d.Spec.Containers,
					util.Volumes:         d.Spec.Volumes,
					util.Labels:          d.ObjectMeta.GetLabels(),
					util.Cluster:         clusterName,
					util.ResourceVersion: d.ObjectMeta.ResourceVersion,
					util.K8sObj:          util.K8sObj,
					util.StartTime:       d.Status.StartTime,
				}
				if len(d.ObjectMeta.OwnerReferences) > 0 {
					pod[util.Owner] = d.ObjectMeta.OwnerReferences[0].Name
					pod[util.OwnerType] = strings.ToLower(d.ObjectMeta.OwnerReferences[0].Kind)
				}
				list = append(list, pod)
			}
			return list, nil
		}
		data := core_v1.Pod{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		pod := map[string]interface{}{
			util.ObjType:         util.Pod,
			util.Name:            data.ObjectMeta.Name,
			util.Namespace:       data.ObjectMeta.Namespace,
			util.CreationTime:    data.ObjectMeta.CreationTimestamp,
			util.Phase:           data.Status.Phase,
			util.NodeName:        data.Spec.NodeName,
			util.IP:              data.Status.PodIP,
			util.Containers:      data.Spec.Containers,
			util.Volumes:         data.Spec.Volumes,
			util.Labels:          data.ObjectMeta.GetLabels(),
			util.Cluster:         clusterName,
			util.ResourceVersion: data.ObjectMeta.ResourceVersion,
			util.K8sObj:          util.K8sObj,
			util.StartTime:       data.Status.StartTime,
		}
		if len(data.ObjectMeta.OwnerReferences) > 0 {
			pod[util.Owner] = data.ObjectMeta.OwnerReferences[0].Name
			pod[util.OwnerType] = strings.ToLower(data.ObjectMeta.OwnerReferences[0].Kind)
		}
		return pod, nil
	case util.ReplicaSet:
		if isArray {
			list := make([]map[string]interface{}, 0)
			data := []v1beta2.ReplicaSet{}
			err := json.Unmarshal(body, &data)
			if err != nil {
				return nil, err
			}
			for _, d := range data {
				replicaset := map[string]interface{}{
					util.ObjType:         util.ReplicaSet,
					util.Name:            d.ObjectMeta.Name,
					util.CreationTime:    d.ObjectMeta.CreationTimestamp,
					util.Namespace:       d.ObjectMeta.Namespace,
					util.NumReplicas:     d.Spec.Replicas,
					util.PodSpec:         d.Spec.Template.Spec,
					util.Owner:           d.ObjectMeta.OwnerReferences[0].Name,
					util.Cluster:         clusterName,
					util.ResourceVersion: d.ObjectMeta.ResourceVersion,
					util.Labels:          d.ObjectMeta.GetLabels(),
					util.K8sObj:          util.K8sObj,
				}
				list = append(list, replicaset)
			}
			return list, nil
		}
		data := v1beta2.ReplicaSet{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			util.ObjType:         util.ReplicaSet,
			util.Name:            data.ObjectMeta.Name,
			util.CreationTime:    data.ObjectMeta.CreationTimestamp,
			util.Namespace:       data.ObjectMeta.Namespace,
			util.NumReplicas:     data.Spec.Replicas,
			util.PodSpec:         data.Spec.Template.Spec,
			util.Owner:           data.ObjectMeta.OwnerReferences[0].Name,
			util.Cluster:         clusterName,
			util.ResourceVersion: data.ObjectMeta.ResourceVersion,
			util.Labels:          data.ObjectMeta.GetLabels(),
			util.K8sObj:          util.K8sObj,
		}, nil
	case util.Service:
		if isArray {
			list := make([]map[string]interface{}, 0)
			data := []core_v1.Service{}
			err := json.Unmarshal(body, &data)
			if err != nil {
				return nil, err
			}
			for _, d := range data {
				service := map[string]interface{}{
					util.ObjType:         util.Service,
					util.Name:            d.ObjectMeta.Name,
					util.Namespace:       d.ObjectMeta.Namespace,
					util.CreationTime:    d.ObjectMeta.CreationTimestamp,
					util.Selector:        d.Spec.Selector,
					util.Labels:          d.ObjectMeta.GetLabels(),
					util.ClusterIP:       d.Spec.ClusterIP,
					util.ServiceType:     d.Spec.Type,
					util.Ports:           d.Spec.Ports,
					util.Cluster:         clusterName,
					util.ResourceVersion: d.ObjectMeta.ResourceVersion,
					util.K8sObj:          util.K8sObj,
				}
				// creata application from labels
				appList := createAppNameList(&d)
				if len(appList) > 0 {
					service[util.Application] = appList
				}
				list = append(list, service)
			}
			return list, nil
		}
		data := core_v1.Service{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		service := map[string]interface{}{
			util.ObjType:         util.Service,
			util.Name:            data.ObjectMeta.Name,
			util.Namespace:       data.ObjectMeta.Namespace,
			util.CreationTime:    data.ObjectMeta.CreationTimestamp,
			util.Selector:        data.Spec.Selector,
			util.Labels:          data.ObjectMeta.GetLabels(),
			util.ClusterIP:       data.Spec.ClusterIP,
			util.ServiceType:     data.Spec.Type,
			util.Ports:           data.Spec.Ports,
			util.Cluster:         clusterName,
			util.ResourceVersion: data.ObjectMeta.ResourceVersion,
			util.K8sObj:          util.K8sObj,
		}
		// creata application from labels
		appList := createAppNameList(&data)
		if len(appList) > 0 {
			service[util.Application] = appList
		}
		return service, nil
	case util.StatefulSet:
		if isArray {
			list := make([]map[string]interface{}, 0)
			data := []appsv1.StatefulSet{}
			err := json.Unmarshal(body, &data)
			if err != nil {
				return nil, err
			}
			for _, d := range data {
				statefulset := map[string]interface{}{
					util.ObjType:         util.StatefulSet,
					util.Name:            d.ObjectMeta.Name,
					util.CreationTime:    d.ObjectMeta.CreationTimestamp,
					util.Namespace:       d.ObjectMeta.Namespace,
					util.NumReplicas:     d.Spec.Replicas,
					util.Cluster:         clusterName,
					util.ResourceVersion: d.ObjectMeta.ResourceVersion,
					util.Labels:          d.ObjectMeta.GetLabels(),
					util.K8sObj:          util.K8sObj,
				}
				list = append(list, statefulset)
			}
			return list, nil
		}
		data := appsv1.StatefulSet{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			util.ObjType:         util.StatefulSet,
			util.Name:            data.ObjectMeta.Name,
			util.CreationTime:    data.ObjectMeta.CreationTimestamp,
			util.Namespace:       data.ObjectMeta.Namespace,
			util.NumReplicas:     data.Spec.Replicas,
			util.Cluster:         clusterName,
			util.ResourceVersion: data.ObjectMeta.ResourceVersion,
			util.Labels:          data.ObjectMeta.GetLabels(),
			util.K8sObj:          util.K8sObj,
		}, nil
	default:
		var data interface{}
		if isArray {
			data = []map[string]interface{}{}
		} else {
			data = map[string]interface{}{}
		}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func getValues(data interface{}, key, method string) string {
	vals := []reflect.Value{}
	switch data.(type) {
	case *core_v1.Service:
		vals = reflect.ValueOf(&data.(*core_v1.Service).ObjectMeta).MethodByName(method).Call(nil)
	case *core_v1.Namespace:
		vals = reflect.ValueOf(&data.(*core_v1.Namespace).ObjectMeta).MethodByName(method).Call(nil)
	case *v1beta2.Deployment:
		vals = reflect.ValueOf(&data.(*v1beta2.Deployment).ObjectMeta).MethodByName(method).Call(nil)
	case *ext_v1beta1.Ingress:
		vals = reflect.ValueOf(&data.(*ext_v1beta1.Ingress).ObjectMeta).MethodByName(method).Call(nil)
	case *core_v1.Pod:
		vals = reflect.ValueOf(&data.(*core_v1.Pod).ObjectMeta).MethodByName(method).Call(nil)
	case *v1beta2.ReplicaSet:
		vals = reflect.ValueOf(&data.(*v1beta2.ReplicaSet).ObjectMeta).MethodByName(method).Call(nil)
	case *appsv1.StatefulSet:
		vals = reflect.ValueOf(&data.(*appsv1.StatefulSet).ObjectMeta).MethodByName(method).Call(nil)
	}
	if len(vals) > 0 {
		if val, ok := vals[0].Interface().(map[string]string)[key]; ok {
			return val
		}
	}
	return ""
}

func createAppNameList(obj interface{}) []interface{} {
	appList := make([]interface{}, 0)
	for _, key := range []string{util.App, util.K8sApp} {
		val := getValues(obj, key, "GetLabels")
		if val != "" {
			appList = append(appList, val)
		}
	}
	return appList
}

// QSLHandler handles requests for QSL
func (s *ServerResource) QSLHandler(w http.ResponseWriter, r *http.Request) {

	metrics.KatlasNumReqCount.Inc()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	// get query for count only
	query, err := s.QSLSvc.CreateDgraphQuery(vars[util.Query], true)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		if err.Error() == "Failed to connect to dgraph to get metadata" {
			metrics.KatlasNumReqErr5xx.Inc()
			w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
			return
		}
		// code: 400
		metrics.KatlasNumReqErr4xx.Inc()
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
		return
	}

	response, err := s.QSLSvc.DBclient.ExecuteDgraphQuery(query)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	total := apis.GetTotalCnt(response)

	// get query with pagination
	query, err = s.QSLSvc.CreateDgraphQuery(vars[util.Query], false)
	log.Infof("dgraph query for %#v:\n %s", vars[util.Query], query)
	if err != nil {
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr4xx.Inc()
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusBadRequest, trim(err.Error()))))
		return
	}
	start := time.Now()
	code := http.StatusOK
	defer func() {
		//metrics.KatlasQueryLatencyHistogram.WithLabelValues(fmt.Sprintf("%d", code)).Observe(time.Since(start).Seconds())
		metrics.KatlasQueryLatencyHistogram.WithLabelValues("katlas", "*", "None", "dev", "containers", "GET", fmt.Sprintf("%d", code), "/**").Observe(time.Since(start).Seconds())
	}()

	response, err = s.QSLSvc.DBclient.ExecuteDgraphQuery(query)
	if err != nil {
		metrics.DgraphNumQSLErr.Inc()
		code = http.StatusInternalServerError
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	log.Infof("[elaps	edtime: %s]response for query %#v", time.Since(start), vars[util.Query])
	response[util.Count] = total
	response["status"] = http.StatusOK
	ret, err := json.Marshal(response)
	if err != nil {
		metrics.DgraphNumQSLErr.Inc()
		code = http.StatusInternalServerError
		metrics.KatlasNumReqErr.Inc()
		metrics.KatlasNumReqErr5xx.Inc()
		w.Write([]byte(fmt.Sprintf("{\"status\": %v, \"error\": \"%s\"}", http.StatusInternalServerError, trim(err.Error()))))
		return
	}
	w.Write(ret)

	metrics.KatlasNumReq2xx.Inc()
}

func trim(str string) string {
	return strings.Replace(strings.Replace(str, "\n", " ", -1), "\"", "'", -1)
}

// TODO:
// Add more supported rest APIs
