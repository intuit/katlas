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

// GetReplicaSetInformer get index Informer to watch ReplicaSet
func GetReplicaSetInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the pods (core resource) in the deafult namespace
				return client.AppsV1beta2().ReplicaSets(AppNamespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the pods (core resource) in the default namespace
				return client.AppsV1beta2().ReplicaSets(AppNamespace).Watch(options)
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

// ValidateReplicaSet check required fields
func ValidateReplicaSet(replicaset *v1beta2.ReplicaSet) bool {
	if replicaset.ObjectMeta.Name == "" {
		return false
	}
	if replicaset.ObjectMeta.Namespace == "" {
		return false
	}
	if replicaset.ObjectMeta.ResourceVersion == "" {
		return false
	}
	return true
}

// ObjectCreated is called when an object is created
func (t *ReplicaSetHandler) ObjectCreated(obj interface{}) error {
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
		return errors.New("Could not validate replicaset object " + replicaset.ObjectMeta.Name)
	}
	SendJSONQueryWithRetries(replicaset, RestSvcEndpoint+"v1/entity/replicaset")
	return nil
}

// ObjectDeleted is called when an object is deleted
func (t *ReplicaSetHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("ReplicaSetHandler.ObjectDeleted")
	SendDeleteRequest(RestSvcEndpoint + "v1/entity/replicaset/replicaset:" + ClusterName + ":" + strings.Replace(key, "/", ":", -1))
	return nil
}

// ObjectUpdated is called when an object is updated
func (t *ReplicaSetHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("ReplicaSetHandler.ObjectUpdated")
	return nil
}

// ReplicaSetSynchronize sync all ReplicaSets periodically in case missing events
func ReplicaSetSynchronize(client kubernetes.Interface) {
	clusterreplicasetslist, _ := client.AppsV1beta2().ReplicaSets(AppNamespace).List(v1.ListOptions{})
	SendJSONQueryWithRetries(clusterreplicasetslist.Items, RestSvcEndpoint+"v1/sync/replicaset")
}
