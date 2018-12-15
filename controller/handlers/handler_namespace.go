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

func ValidateNamespace(pod *core_v1.Namespace) bool {
	if pod.ObjectMeta.Name == "" {
		return false
	}

	if pod.ObjectMeta.ResourceVersion == "" {
		return false
	}

	return true
}

// ObjectCreated is called when an object is created
func (t *NamespaceHandler) ObjectCreated(obj interface{}) (map[string]interface{}, error) {
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
		return nil, errors.New("Could not validate namespace object " + namespace.ObjectMeta.Name)
	}
	namespacedata := CreateNamespaceData(*namespace, CLUSTERNAME)

	SendJSONQueryWithRetries(namespacedata, CMDBAPIENDPOINT+"v1/entity/Namespace")

	return namespacedata, nil
}

// ObjectDeleted is called when an object is deleted
func (t *NamespaceHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("NamespaceHandler.ObjectDeleted")

	SendDeleteRequest(CMDBAPIENDPOINT + "v1/entity/Namespace/" + CLUSTERNAME + ":" + key)

	return nil
}

// ObjectUpdated is called when an object is updated
func (t *NamespaceHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("NamespaceHandler.ObjectUpdated")
	return nil
}

func CreateNamespaceData(namespace core_v1.Namespace, clustername string) map[string]interface{} {

	namespacedata := map[string]interface{}{
		"objtype":         "Namespace",
		"name":            namespace.ObjectMeta.Name,
		"creationtime":    namespace.ObjectMeta.CreationTimestamp,
		"cluster":         clustername,
		"resourceversion": namespace.ResourceVersion,
		"k8sobj":          "K8sObj",
		"labels":          namespace.ObjectMeta.GetLabels(),
	}

	return namespacedata
}

func NamespaceSynchronize(client kubernetes.Interface) {
	list := make([]map[string]interface{}, 0)
	clusternamespaceslist, _ := client.CoreV1().Namespaces().List(v1.ListOptions{})
	for _, data := range clusternamespaceslist.Items {
		list = append(list, CreateNamespaceData(data, CLUSTERNAME))
	}
	SendJSONQueryWithRetries(list, CMDBAPIENDPOINT+"v1/sync/Namespace")
}
