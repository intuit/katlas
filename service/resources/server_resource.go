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
	clusterName := r.Header.Get(util.ClusterName)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	payload, err := buildEntityData(clusterName, meta, body, false)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uids, err := s.EntitySvc.CreateEntity(meta, payload.(map[string]interface{}))
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

// EntitySyncHandler REST API to sync entities
func (s ServerResource) EntitySyncHandler(w http.ResponseWriter, r *http.Request) {
	//Set Access-Control-Allow-Origin header now so that it will be present
	//even if an error is returned (otherwise the error also causes a CORS
	//exception in the browser/client)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	meta := vars[util.Metadata]
	clusterName := r.Header.Get(util.ClusterName)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	payload, err := buildEntityData(clusterName, meta, body, true)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.EntitySvc.SyncEntities(meta, payload.([]map[string]interface{}))
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	var payload interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if reflect.TypeOf(payload).Kind() == reflect.Slice {
		var rets []string
		for _, p := range payload.([]interface{}) {
			_, err := s.MetaSvc.CreateMetadata(p.(map[string]interface{}))
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			rets = append(rets, p.(map[string]interface{})[util.Name].(string))
		}
		w.Write([]byte(fmt.Sprintf("Metadata %v create successfully", rets)))
	} else {
		_, err := s.MetaSvc.CreateMetadata(payload.(map[string]interface{}))
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("Metadata %s create successfully", payload.(map[string]interface{})[util.Name])))
	}
}

// SchemaCreateHandler REST API for create Schema
func (s ServerResource) SchemaCreateHandler(w http.ResponseWriter, r *http.Request) {
	defer s.MetaSvc.RemoveSchemaCache(db.LruCache)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	var payload interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if reflect.TypeOf(payload).Kind() == reflect.Slice {
		var predicates []db.Schema
		err := mapstructure.Decode(payload, &predicates)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, p := range predicates {
			err := s.MetaSvc.CreateSchema(p)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		var predicate db.Schema
		err := mapstructure.Decode(payload, &predicate)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = s.MetaSvc.CreateSchema(predicate)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("Schema create successfully"))
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
		}, nil
	case util.Deployment:
		if isArray {
			list := make([]map[string]interface{}, 0)
			data := []ext_v1beta1.Deployment{}
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

func createAppNameList(obj interface{}) []interface{} {
	appList := make([]interface{}, 0)
	switch obj.(type) {
	case *core_v1.Service:
		if appName, ok := obj.(*core_v1.Service).ObjectMeta.GetLabels()[util.App]; ok {
			appList = append(appList, appName)
		}
		if appName, ok := obj.(*core_v1.Service).ObjectMeta.GetLabels()[util.K8sApp]; ok {
			appList = append(appList, appName)
		}
	case *ext_v1beta1.Ingress:
		if appName, ok := obj.(*ext_v1beta1.Ingress).ObjectMeta.GetLabels()[util.App]; ok {
			appList = append(appList, appName)
		}
		if appName, ok := obj.(*ext_v1beta1.Ingress).ObjectMeta.GetLabels()[util.K8sApp]; ok {
			appList = append(appList, appName)
		}
	case *v1beta2.Deployment:
		if appName, ok := obj.(*v1beta2.Deployment).ObjectMeta.GetLabels()[util.App]; ok {
			appList = append(appList, appName)
		}
		if appName, ok := obj.(*v1beta2.Deployment).ObjectMeta.GetLabels()[util.K8sApp]; ok {
			appList = append(appList, appName)
		}
	}
	return appList
}

// QSLHandler handles requests for QSL
func (s *ServerResource) QSLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)

	// get query for count only
	query, err := s.QSLSvc.CreateDgraphQuery(vars[util.Query], true)
	if err != nil {
		if err.Error() == "Failed to connect to dgraph to get metadata" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest) // code: 400
		return
	}

	response, err := s.QSLSvc.DBclient.ExecuteDgraphQuery(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var total float64
	for _, res := range response[util.Objects].([]interface{}) {
		val, ok := res.(map[string]interface{})[util.Count]
		if ok {
			total = val.(float64)
		}
	}

	// get query with pagination
	query, err = s.QSLSvc.CreateDgraphQuery(vars[util.Query], false)
	log.Infof("dgraph query for %#v:\n %s", vars[util.Query], query)
	start := time.Now()
	response, err = s.QSLSvc.DBclient.ExecuteDgraphQuery(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("[elapsedtime: %s]response for query %#v", time.Since(start), vars[util.Query])
	response[util.Count] = total
	ret, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to convert to JSON output", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

// TODO:
// Add more supported rest APIs
