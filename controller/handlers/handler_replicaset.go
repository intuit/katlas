package handlers

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	"k8s.io/api/apps/v1beta2"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	//metav1 "k8s.io/apimachinery/pkg/resources/meta/v1"
)

// ReplicaSetHandler is a sample implementation of Handler
type ReplicaSetHandler struct{}

func GetReplicaSetInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the pods (core resource) in the deafult namespace
				return client.AppsV1beta2().ReplicaSets(APP_NAMESPACE).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the pods (core resource) in the default namespace
				return client.AppsV1beta2().ReplicaSets(APP_NAMESPACE).Watch(options)
			},
		},
		&v1beta2.ReplicaSet{}, // the target type (Pod)
		0,                     // no resync (period of 0)
		cache.Indexers{},
	)
	return informer
}

// Init handles any handler initialization
// a method of ReplicaSetHandler returns type error
// func (<object>) <name>(<params>) <return>
func (t *ReplicaSetHandler) Init() error {
	log.Info("ReplicaSetHandler.Init")
	return nil
}

func ValidateReplicaSet(pod *v1beta2.ReplicaSet) bool {
	if pod.ObjectMeta.Name == "" {
		return false
	}

	if pod.ObjectMeta.Namespace == "" {
		return false
	}

	if pod.ObjectMeta.ResourceVersion == "" {
		return false
	}

	return true
}

// ObjectCreated is called when an object is created
func (t *ReplicaSetHandler) ObjectCreated(obj interface{}) (map[string]interface{}, error) {
	log.Info("ReplicaSetHandler.ObjectCreated")
	defer func() {
		if r := recover(); r != nil {
			t.ObjectUpdated(obj, obj)
			return
		}
	}()
	// assert the type to a Pod object to pull out relevant data
	replicaset := obj.(*v1beta2.ReplicaSet)

	if !ValidateReplicaSet(replicaset) {
		return nil, errors.New("Could not validate replicaset object " + replicaset.ObjectMeta.Name)
	}

	replicasetdata := CreateReplicaSetData(*replicaset, CLUSTERNAME)

	SendJSONQueryWithRetries(replicasetdata, CMDBAPIENDPOINT+"v1/entity/ReplicaSet")

	return replicasetdata, nil
}

// ObjectDeleted is called when an object is deleted
func (t *ReplicaSetHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("ReplicaSetHandler.ObjectDeleted")

	SendDeleteRequest(CMDBAPIENDPOINT + "v1/entity/ReplicaSet/" + CLUSTERNAME + ":" + strings.Replace(key, "/", ":", -1))

	return nil
}

// ObjectUpdated is called when an object is updated
func (t *ReplicaSetHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("ReplicaSetHandler.ObjectUpdated")
	return nil
}

func CreateReplicaSetData(replicaset v1beta2.ReplicaSet, clustername string) map[string]interface{} {

	replicasetdata := map[string]interface{}{
		"objtype":         "ReplicaSet",
		"name":            replicaset.ObjectMeta.Name,
		"creationtime":    &replicaset.ObjectMeta.CreationTimestamp,
		"namespace":       replicaset.ObjectMeta.Namespace,
		"numreplicas":     replicaset.Spec.Replicas,
		"podspec":         replicaset.Spec.Template.Spec,
		"owner":           replicaset.ObjectMeta.OwnerReferences[0].Name,
		"cluster":         clustername,
		"resourceversion": replicaset.ObjectMeta.ResourceVersion,
		"labels":          replicaset.ObjectMeta.GetLabels(),
		"k8sobj":          "K8sObj",
	}

	return replicasetdata
}

func ReplicaSetSynchronize(client kubernetes.Interface) {
	list := make([]map[string]interface{}, 0)
	clusterreplicasetslist, _ := client.AppsV1beta2().ReplicaSets(APP_NAMESPACE).List(v1.ListOptions{})
	for _, data := range clusterreplicasetslist.Items {
		list = append(list, CreateReplicaSetData(data, CLUSTERNAME))
	}
	SendJSONQueryWithRetries(list, CMDBAPIENDPOINT+"v1/sync/ReplicaSet")
}
