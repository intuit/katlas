package handlers

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	//metav1 "k8s.io/apimachinery/pkg/resources/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// IngressHandler is a sample implementation of Handler
type IngressHandler struct{}

// GetIngressInformer get index Informer to watch Ingress
func GetIngressInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the ingresses (ExtensionsV1beta1 resource) in the deafult ingress
				return client.ExtensionsV1beta1().Ingresses(AppNamespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the ingresses (ExtensionsV1beta1 resource) in the default ingress
				return client.ExtensionsV1beta1().Ingresses(AppNamespace).Watch(options)
			},
		},
		&ext_v1beta1.Ingress{}, // the target type (Pod)
		0,                      // no resync (period of 0)
		cache.Indexers{},
	)
	return informer
}

// Init handles any handler initialization
// a method of IngressHandler returns type error
// func (<object>) <name>(<params>) <return>
func (t *IngressHandler) Init() error {
	log.Info("IngressHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *IngressHandler) ObjectCreated(obj interface{}) error {
	log.Info("IngressHandler.ObjectCreated")
	defer func() {
		if r := recover(); r != nil {
			t.ObjectUpdated(obj, obj)
			return
		}
	}()
	// assert the type to a Ingress object to pull out relevant data
	ingress := obj.(*ext_v1beta1.Ingress)
	log.Infof("    ingressmeta: %+v", ingress)
	SendJSONQueryWithRetries(ingress, RestSvcEndpoint+"v1/entity/ingress")
	return nil
}

// ObjectDeleted is called when an object is deleted
func (t *IngressHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("IngressHandler.ObjectDeleted")
	SendDeleteRequest(RestSvcEndpoint + "v1/entity/ingress/ingress:" + ClusterName + ":" + strings.Replace(key, "/", ":", -1))
	return nil
}

// ObjectUpdated is called when an object is updated
func (t *IngressHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("IngressHandler.ObjectUpdated")
	return nil
}

// IngressSynchronize sync all Ingresses periodically in case missing events
func IngressSynchronize(client kubernetes.Interface) {
	clusteringresseslist, _ := client.ExtensionsV1beta1().Ingresses(AppNamespace).List(v1.ListOptions{})
	SendJSONQueryWithRetries(clusteringresseslist.Items, RestSvcEndpoint+"v1/sync/ingress")
}
