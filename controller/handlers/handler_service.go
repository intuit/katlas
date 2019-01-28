package handlers

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/resources/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// ServiceHandler is a sample implementation of Handler
type ServiceHandler struct{}

// GetServiceInformer get index Informer to watch Service
func GetServiceInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the pods (core resource) in the deafult namespace
				return client.CoreV1().Services(AppNamespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the pods (core resource) in the default namespace
				return client.CoreV1().Services(AppNamespace).Watch(options)
			},
		},
		&core_v1.Service{}, // the target type (Pod)
		0,                  // no resync (period of 0)
		cache.Indexers{},
	)

	return informer
}

// Init handles any handler initialization
// a method of ServiceHandler returns type error
// func (<object>) <name>(<params>) <return>
func (t *ServiceHandler) Init() error {
	log.Info("ServiceHandler.Init")
	return nil
}

// ValidateService check required fields
func ValidateService(service *core_v1.Service) bool {
	if service.ObjectMeta.Name == "" {
		return false
	}
	if service.ObjectMeta.Namespace == "" {
		return false
	}
	if service.ObjectMeta.ResourceVersion == "" {
		return false
	}
	return true
}

// ObjectCreated is called when an object is created
func (t *ServiceHandler) ObjectCreated(obj interface{}) error {
	log.Info("ServiceHandler.ObjectCreated")
	defer func() {
		if r := recover(); r != nil {
			t.ObjectUpdated(obj, obj)
			return
		}
	}()
	// assert the type to a Service object to pull out relevant data
	service := obj.(*core_v1.Service)
	if !ValidateService(service) {
		return errors.New("Could not validate service object " + service.ObjectMeta.Name)
	}
	SendJSONQueryWithRetries(service, RestSvcEndpoint+"v1/entity/service")
	return nil
}

// ObjectDeleted is called when an object is deleted
func (t *ServiceHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("ServiceHandler.ObjectDeleted")
	SendDeleteRequest(RestSvcEndpoint + "v1/entity/service/service:" + ClusterName + ":" + strings.Replace(key, "/", ":", -1))
	return nil
}

// ObjectUpdated is called when an object is updated
func (t *ServiceHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("ServiceHandler.ObjectUpdated")
	return nil
}

// ServiceSynchronize sync all Services periodically in case missing events
func ServiceSynchronize(client kubernetes.Interface) {
	clusterserviceslist, _ := client.CoreV1().Services(AppNamespace).List(v1.ListOptions{})
	SendJSONQueryWithRetries(clusterserviceslist.Items, RestSvcEndpoint+"v1/sync/service")
}
