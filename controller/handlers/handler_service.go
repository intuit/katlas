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

func GetServiceInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the pods (core resource) in the deafult namespace
				return client.CoreV1().Services(APP_NAMESPACE).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the pods (core resource) in the default namespace
				return client.CoreV1().Services(APP_NAMESPACE).Watch(options)
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

func ValidateService(pod *core_v1.Service) bool {
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
func (t *ServiceHandler) ObjectCreated(obj interface{}) (map[string]interface{}, error) {
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
		return nil, errors.New("Could not validate service object " + service.ObjectMeta.Name)
	}

	servicedata := CreateServiceData(*service, CLUSTERNAME)

	SendJSONQueryWithRetries(servicedata, CMDBAPIENDPOINT+"v1/entity/Service")

	return servicedata, nil
}

// ObjectDeleted is called when an object is deleted
func (t *ServiceHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("ServiceHandler.ObjectDeleted")

	SendDeleteRequest(CMDBAPIENDPOINT + "v1/entity/Service/" + CLUSTERNAME + ":" + strings.Replace(key, "/", ":", -1))

	return nil
}

// ObjectUpdated is called when an object is updated
func (t *ServiceHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("ServiceHandler.ObjectUpdated")
	return nil
}

func CreateServiceData(service core_v1.Service, clustername string) map[string]interface{} {

	servicedata := map[string]interface{}{
		"objtype":         "Service",
		"name":            service.ObjectMeta.Name,
		"namespace":       service.ObjectMeta.Namespace,
		"creationtime":    service.ObjectMeta.CreationTimestamp,
		"selector":        service.Spec.Selector,
		"labels":          service.ObjectMeta.GetLabels(),
		"clusterip":       service.Spec.ClusterIP,
		"servicetype":     service.Spec.Type,
		"ports":           service.Spec.Ports,
		"cluster":         clustername,
		"resourceversion": service.ObjectMeta.ResourceVersion,
		"k8sobj":          "K8sObj",
	}

	return servicedata
}

func ServiceSynchronize(client kubernetes.Interface) {
	list := make([]map[string]interface{}, 0)
	clusterserviceslist, _ := client.CoreV1().Services(APP_NAMESPACE).List(v1.ListOptions{})
	for _, data := range clusterserviceslist.Items {
		list = append(list, CreateServiceData(data, CLUSTERNAME))
	}
	SendJSONQueryWithRetries(list, CMDBAPIENDPOINT+"v1/sync/Service")
}
