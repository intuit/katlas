package handlers

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	//"reflect"
	//"github.com/dgraph-io/dgo/y"
)

// PodHandler is a sample implementation of Handler
type PodHandler struct{}

// GetPodInformer get index Informer to watch Pod
func GetPodInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	// the informer is responsible for listing and watching events for objects of a specific type
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the pods (core resource) in the deafult namespace
				return client.CoreV1().Pods(AppNamespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the pods (core resource) in the default namespace
				return client.CoreV1().Pods(AppNamespace).Watch(options)
			},
		},
		&core_v1.Pod{}, // the target type (Pod)
		0,              // no resync (period of 0)
		cache.Indexers{},
	)
	return informer
}

// Init handles any handler initialization
// a method of PodHandler returns type error
// func (<object>) <name>(<params>) <return>
func (t *PodHandler) Init() error {
	log.Info("PodHandler.Init")
	return nil
}

// ValidatePod verify that the object has at least these fields
func ValidatePod(pod *core_v1.Pod) bool {
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
func (t *PodHandler) ObjectCreated(obj interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			t.ObjectUpdated(obj, obj)
			return
		}
	}()
	// assert the type to a Pod object to pull out relevant data
	pod := obj.(*core_v1.Pod)
	if !ValidatePod(pod) {
		return errors.New("Could not validate pod object " + pod.ObjectMeta.Name)
	}
	// send the object to the rest service
	SendJSONQueryWithRetries(pod, RestSvcEndpoint+"v1/entity/pod")
	return nil
}

// ObjectDeleted is called when an object is deleted
func (t *PodHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("PodHandler.ObjectDeleted")
	SendDeleteRequest(RestSvcEndpoint + "v1/entity/pod/pod:" + ClusterName + ":" + strings.Replace(key, "/", ":", -1))
	return nil
}

// ObjectUpdated is called when an object is updated
func (t *PodHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("PodHandler.ObjectUpdated")
	return nil
}

// PodSynchronize synchronize the objects in dgraph with the cluster to account for drift
// e.g. if there were network issues and some events weren't received,
// or if the api crashes while processing some events
func PodSynchronize(client kubernetes.Interface) {
	clusterpodslist, _ := client.CoreV1().Pods(AppNamespace).List(v1.ListOptions{})
	SendJSONQueryWithRetries(clusterpodslist.Items, RestSvcEndpoint+"v1/sync/pod")
}
