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

func GetPodInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	// the informer is responsible for listing and watching events for objects of a specific type
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the pods (core resource) in the deafult namespace
				return client.CoreV1().Pods(APP_NAMESPACE).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the pods (core resource) in the default namespace
				return client.CoreV1().Pods(APP_NAMESPACE).Watch(options)
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

// verify that the object has at least these fields
func ValidatePod(pod *core_v1.Pod) bool {
	if pod.ObjectMeta.Name == "" {
		return false
	}
	if pod.ObjectMeta.Namespace == "" {
		return false
	}
	if len(pod.ObjectMeta.OwnerReferences) == 0 {
		return false
	}

	return true
}

// ObjectCreated is called when an object is created
func (t *PodHandler) ObjectCreated(obj interface{}) (map[string]interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			t.ObjectUpdated(obj, obj)
			return
		}
	}()
	// assert the type to a Pod object to pull out relevant data
	pod := obj.(*core_v1.Pod)
	if !ValidatePod(pod) {
		return nil, errors.New("Could not validate pod object " + pod.ObjectMeta.Name)
	}

	poddata := CreatePodData(*pod, CLUSTERNAME)

	// send the object to the rest service
	SendJSONQueryWithRetries(poddata, CMDBAPIENDPOINT+"v1/entity/Pod")

	return poddata, nil
}

// ObjectDeleted is called when an object is deleted
func (t *PodHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("PodHandler.ObjectDeleted")

	SendDeleteRequest(CMDBAPIENDPOINT + "v1/entity/Pod/" + CLUSTERNAME + ":" + strings.Replace(key, "/", ":", -1))

	return nil
}

// ObjectUpdated is called when an object is updated
func (t *PodHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("PodHandler.ObjectUpdated")
	return nil
}

func CreatePodData(pod core_v1.Pod, clustername string) map[string]interface{} {

	podmap := map[string]interface{}{
		"objtype":         "Pod",
		"name":            pod.ObjectMeta.Name,
		"namespace":       pod.ObjectMeta.Namespace,
		"creationtime":    pod.ObjectMeta.CreationTimestamp,
		"phase":           pod.Status.Phase,
		"nodename":        pod.Spec.NodeName,
		"ip":              pod.Status.PodIP,
		"containers":      pod.Spec.Containers,
		"volumes":         pod.Spec.Volumes,
		"labels":          pod.ObjectMeta.GetLabels(),
		"cluster":         clustername,
		"resourceversion": pod.ObjectMeta.ResourceVersion,
		"k8sobj":          "K8sObj",
	}
	if len(pod.ObjectMeta.OwnerReferences) > 0 {
		podmap["owner"] = pod.ObjectMeta.OwnerReferences[0].Name
		podmap["ownertype"] = pod.ObjectMeta.OwnerReferences[0].Kind
	}

	return podmap
}

// synchronize the objects in dgraph with the cluster to account for drift
// e.g. if there were network issues and some events weren't received,
// or if the api crashes while processing some events
func PodSynchronize(client kubernetes.Interface) {
	list := make([]map[string]interface{}, 0)
	clusterpodslist, _ := client.CoreV1().Pods(APP_NAMESPACE).List(v1.ListOptions{})
	for _, data := range clusterpodslist.Items {
		list = append(list, CreatePodData(data, CLUSTERNAME))
	}
	SendJSONQueryWithRetries(list, CMDBAPIENDPOINT+"v1/sync/Pod")
}
