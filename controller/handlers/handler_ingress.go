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

func GetIngressInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the ingresses (ExtensionsV1beta1 resource) in the deafult ingress
				return client.ExtensionsV1beta1().Ingresses(APP_NAMESPACE).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the ingresses (ExtensionsV1beta1 resource) in the default ingress
				return client.ExtensionsV1beta1().Ingresses(APP_NAMESPACE).Watch(options)
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
func (t *IngressHandler) ObjectCreated(obj interface{}) (map[string]interface{}, error) {
	log.Info("IngressHandler.ObjectCreated")
	defer func() {
		if r := recover(); r != nil {
			t.ObjectUpdated(obj, obj)
			return
		}
	}()
	// assert the type to a Ingress object to pull out relevant data
	ingress := obj.(*ext_v1beta1.Ingress)

	ingressdata := CreateIngressData(*ingress, CLUSTERNAME)

	log.Infof("    ingressmeta: %+v", ingress)
	log.Infof(".   IngressName: %s", ingressdata)

	SendJSONQueryWithRetries(ingressdata, CMDBAPIENDPOINT+"v1/entity/Ingress")

	return ingressdata, nil
}

// ObjectDeleted is called when an object is deleted
func (t *IngressHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("IngressHandler.ObjectDeleted")

	SendDeleteRequest(CMDBAPIENDPOINT + "v1/entity/Ingress/" + CLUSTERNAME + ":" + strings.Replace(key, "/", ":", -1))

	return nil
}

// ObjectUpdated is called when an object is updated
func (t *IngressHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("IngressHandler.ObjectUpdated")
	return nil
}

func CreateIngressData(ingress ext_v1beta1.Ingress, clustername string) map[string]interface{} {

	ingressdata := map[string]interface{}{
		"objtype":         "Ingress",
		"cluster":         clustername,
		"name":            ingress.ObjectMeta.Name,
		"namespace":       ingress.ObjectMeta.Namespace,
		"creationtime":    ingress.ObjectMeta.CreationTimestamp,
		"defaultbackend":  ingress.Spec.Backend,
		"tls":             ingress.Spec.TLS,
		"rules":           ingress.Spec.Rules,
		"resourceversion": ingress.ObjectMeta.ResourceVersion,
		"labels":          ingress.ObjectMeta.GetLabels(),
		"k8sobj":          "K8sObj",
	}

	return ingressdata
}

func IngressSynchronize(client kubernetes.Interface) {
	list := make([]map[string]interface{}, 0)
	clusteringresseslist, _ := client.ExtensionsV1beta1().Ingresses(APP_NAMESPACE).List(v1.ListOptions{})
	for _, data := range clusteringresseslist.Items {
		list = append(list, CreateIngressData(data, CLUSTERNAME))
	}
	SendJSONQueryWithRetries(list, CMDBAPIENDPOINT+"v1/sync/Ingress")
}
