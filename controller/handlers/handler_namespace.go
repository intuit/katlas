package handlers

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/resources/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// NamespaceHandler is a sample implementation of Handler
type NamespaceHandler struct{}

// GetNamespaceInformer get index Informer to watch Namespace
func GetNamespaceInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the namespaces (core resource) in the deafult namespace
				return client.CoreV1().Namespaces().List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the namespaces (core resource) in the default namespace
				return client.CoreV1().Namespaces().Watch(options)
			},
		},
		&core_v1.Namespace{}, // the target type (Pod)
		0,                    // no resync (period of 0)
		cache.Indexers{},
	)
	return informer
}

// Init handles any handler initialization
// a method of NamespaceHandler returns type error
// func (<object>) <name>(<params>) <return>
func (t *NamespaceHandler) Init() error {
	log.Info("NamespaceHandler.Init")
	return nil
}

// ValidateNamespace check required attributes
func ValidateNamespace(namespace *core_v1.Namespace) bool {
	if namespace.ObjectMeta.Name == "" {
		return false
	}
	if namespace.ObjectMeta.ResourceVersion == "" {
		return false
	}
	return true
}

// ObjectCreated is called when an object is created
func (t *NamespaceHandler) ObjectCreated(obj interface{}) error {
	log.Info("NamespaceHandler.ObjectCreated")
	defer func() {
		if r := recover(); r != nil {
			t.ObjectUpdated(obj, obj)
			return
		}
	}()
	// assert the type to a Namespace object to pull out relevant data
	namespace := obj.(*core_v1.Namespace)
	if !ValidateNamespace(namespace) {
		return errors.New("Could not validate namespace object " + namespace.ObjectMeta.Name)
	}
	SendJSONQueryWithRetries(namespace, RestSvcEndpoint+"v1/entity/namespace")
	return nil
}

// ObjectDeleted is called when an object is deleted
func (t *NamespaceHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("NamespaceHandler.ObjectDeleted")
	SendDeleteRequest(RestSvcEndpoint + "v1/entity/namespace/namespace:" + ClusterName + ":" + key)
	return nil
}

// ObjectUpdated is called when an object is updated
func (t *NamespaceHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("NamespaceHandler.ObjectUpdated")
	return nil
}

// NamespaceSynchronize sync all Namespaces periodically in case missing events
func NamespaceSynchronize(client kubernetes.Interface) {
	clusternamespaceslist, _ := client.CoreV1().Namespaces().List(v1.ListOptions{})
	SendJSONQueryWithRetries(clusternamespaceslist.Items, RestSvcEndpoint+"v1/sync/namespace")
}
