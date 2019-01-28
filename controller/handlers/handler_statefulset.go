package handlers

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	//metav1 "k8s.io/apimachinery/pkg/resources/meta/v1"
)

// StatefulSetHandler is a sample implementation of Handler
type StatefulSetHandler struct{}

// GetStatefulSetInformer get index Informer to watch StatefulSet
func GetStatefulSetInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the statefulsets (apps resource) in the deafult namespace
				return client.AppsV1().StatefulSets(AppNamespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the statefulsets (apps resource) in the default namespace
				return client.AppsV1().StatefulSets(AppNamespace).Watch(options)
			},
		},
		&appsv1.StatefulSet{}, // the target type (Pod)
		0,                     // no resync (period of 0)
		cache.Indexers{},
	)
	return informer
}

// Init handles any handler initialization
// a method of StatefulSetHandler returns type error
// func (<object>) <name>(<params>) <return>
func (t *StatefulSetHandler) Init() error {
	log.Info("StatefulSetHandler.Init")
	return nil
}

// ValidateStatefulSet to check required fields
func ValidateStatefulSet(sts *appsv1.StatefulSet) bool {
	if sts.ObjectMeta.Name == "" {
		return false
	}
	if sts.ObjectMeta.Namespace == "" {
		return false
	}
	if sts.ObjectMeta.ResourceVersion == "" {
		return false
	}
	return true
}

// ObjectCreated is called when an object is created
func (t *StatefulSetHandler) ObjectCreated(obj interface{}) error {
	log.Info("StatefulSetHandler.ObjectCreated")
	defer func() {
		if r := recover(); r != nil {
			t.ObjectUpdated(obj, obj)
			return
		}
	}()
	// assert the type to a StatefulSet object to pull out relevant data
	statefulset := obj.(*appsv1.StatefulSet)

	if !ValidateStatefulSet(statefulset) {
		return errors.New("Could not validate statefulset object " + statefulset.ObjectMeta.Name)
	}
	SendJSONQueryWithRetries(statefulset, RestSvcEndpoint+"v1/entity/statefulset")
	return nil
}

// ObjectDeleted is called when an object is deleted
func (t *StatefulSetHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("StatefulSetHandler.ObjectDeleted")
	SendDeleteRequest(RestSvcEndpoint + "v1/entity/statefulset/statefulset:" + ClusterName + ":" + strings.Replace(key, "/", ":", -1))
	return nil
}

// ObjectUpdated is called when an object is updated
func (t *StatefulSetHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("StatefulSetHandler.ObjectUpdated")
	return nil
}

// StatefulSetSynchronize sync all StatefulSets periodically in case missing events
func StatefulSetSynchronize(client kubernetes.Interface) {
	clusterstatefulsetslist, _ := client.AppsV1().StatefulSets(AppNamespace).List(v1.ListOptions{})
	SendJSONQueryWithRetries(clusterstatefulsetslist.Items, RestSvcEndpoint+"v1/sync/statefulset")
}
