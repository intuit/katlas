package handlers

import (
	"encoding/json"
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	v1beta2 "k8s.io/api/apps/v1beta2"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
	//metav1 "k8s.io/apimachinery/pkg/resources/meta/v1"
)

// DeploymentHandler is a sample implementation of Handler
type DeploymentHandler struct{}

// GetDeploymentInformer get index Informer to watch Deployment
func GetDeploymentInformer(client kubernetes.Interface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		// the ListWatch contains two different functions that our
		// informer requires: ListFunc to take care of listing and watching
		// the resources we want to handle
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				// list all of the deployments (AppsV1beta2 resource) in the deafult namespace
				return client.AppsV1beta2().Deployments(AppNamespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				// watch all of the deployments (AppsV1beta2 resource) in the default namespace
				return client.AppsV1beta2().Deployments(AppNamespace).Watch(options)
			},
		},
		&v1beta2.Deployment{}, // the target type (Pod)
		0,                     // no resync (period of 0)
		cache.Indexers{},
	)
	return informer
}

// Init handles any handler initialization
// a method of DeploymentHandler returns type error
// func (<object>) <name>(<params>) <return>
func (t *DeploymentHandler) Init() error {
	log.Info("DeploymentHandler.Init")
	return nil
}

// ValidateDeployment to check deployment required attributes
func ValidateDeployment(deployment *v1beta2.Deployment) bool {
	if deployment.ObjectMeta.Name == "" {
		return false
	}
	if deployment.ObjectMeta.Namespace == "" {
		return false
	}
	if deployment.ResourceVersion == "" {
		return false
	}
	return true
}

// ObjectCreated is called when an object is created
func (t *DeploymentHandler) ObjectCreated(obj interface{}) error {
	log.Info("DeploymentHandler.ObjectCreated")
	// assert the type to a Pod object to pull out relevant data
	deployment := obj.(*v1beta2.Deployment)
	defer func() {
		if r := recover(); r != nil {
			t.ObjectUpdated(obj, obj)
			return
		}
	}()
	if !ValidateDeployment(deployment) {
		return errors.New("Could not validate deployment object " + deployment.ObjectMeta.Name)
	}
	j, err := json.MarshalIndent(deployment, "", "    ")
	if err != nil {
		log.Error(err)
	}
	log.Debugf("    Deployment: %s, \n", j)
	SendJSONQueryWithRetries(deployment, RestSvcEndpoint+"v1/entity/deployment")
	return nil
}

// ObjectDeleted is called when an object is deleted
func (t *DeploymentHandler) ObjectDeleted(obj interface{}, key string) error {
	log.Info("DeploymentHandler.ObjectDeleted")
	SendDeleteRequest(RestSvcEndpoint + "v1/entity/deployment/deployment:" + ClusterName + ":" + strings.Replace(key, "/", ":", -1))
	return nil
}

// ObjectUpdated is called when an object is updated
func (t *DeploymentHandler) ObjectUpdated(objOld, objNew interface{}) error {
	log.Info("DeploymentHandler.ObjectUpdated")
	return nil
}

// DeploymentSynchronize sync all Deployments periodically in case missing events
func DeploymentSynchronize(client kubernetes.Interface) {
	clusterdeploymentslist, _ := client.AppsV1beta2().Deployments(AppNamespace).List(v1.ListOptions{})
	SendJSONQueryWithRetries(clusterdeploymentslist.Items, RestSvcEndpoint+"v1/sync/deployment")
}
