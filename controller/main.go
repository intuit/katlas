package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	handlers "github.com/intuit/katlas/controller/handlers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

// retrieve the Kubernetes cluster client from outside of the cluster
func GetKubernetesClient() kubernetes.Interface {
	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"
	// kubeConfigPath := os.Getenv("HOME") + "/Downloads/admins\\@dev-devx-cmdb-api-usw2-ppd-qal"
	var config *rest.Config
	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Infof("getClusterConfig: %v", err)

		// if the container is running inside the cluster, use the incluster config
		var err2 error
		config, err2 = rest.InClusterConfig()
		if err2 != nil {
			panic(err.Error())
		}
	}

	// generate the client based off of the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client
}

func CreateController(obj_type string) *Controller {
	client := GetKubernetesClient()

	var informer cache.SharedIndexInformer
	var handlerc handlers.Handler

	// create a new queue so that when the informer gets a resource that is either
	// a result of listing or watching, we can add an idenfitying key to the queue
	// so that it can be handled in the handler
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	switch obj_type {
	case "Pod":
		// create the informer so that we can not only list resources
		// but also watch them for all pods in the default namespace
		informer = handlers.GetPodInformer(client)

		// add event handlers to handle the three types of events for resources:
		//  - adding new resources
		//  - updating existing resources
		//  - deleting resources
		handlerc = &handlers.PodHandler{}
		//return informer
	case "Service":
		informer = handlers.GetServiceInformer(client)
		handlerc = &handlers.ServiceHandler{}

	case "Namespace":
		informer = handlers.GetNamespaceInformer(client)
		handlerc = &handlers.NamespaceHandler{}

	case "Deployment":
		informer = handlers.GetDeploymentInformer(client)
		handlerc = &handlers.DeploymentHandler{}

	case "ReplicaSet":
		informer = handlers.GetReplicaSetInformer(client)
		handlerc = &handlers.ReplicaSetHandler{}

	case "Ingress":
		informer = handlers.GetIngressInformer(client)
		handlerc = &handlers.IngressHandler{}

	case "StatefulSet":
		informer = handlers.GetStatefulSetInformer(client)
		handlerc = &handlers.StatefulSetHandler{}
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// convert the resource object into a key (in this case
			// we are just doing it in the format of 'namespace/name')
			key, err := cache.MetaNamespaceKeyFunc(obj)
			log.Infof("Add %s: %s", obj_type, key)
			if err == nil {
				// add the key to the queue for the handler to get
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			log.Infof("Update %s: %s", obj_type, key)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// DeletionHandlingMetaNamsespaceKeyFunc is a helper function that allows
			// us to check the DeletedFinalStateUnknown existence in the event that
			// a resource was deleted but it is still contained in the index
			//
			// this then in turn calls MetaNamespaceKeyFunc
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			log.Infof("Delete %s: %s", obj_type, key)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	// construct the Controller object which has all of the necessary components to
	// handle logging, connections, informing (listing and watching), the queue,
	// and the handler
	controller := Controller{
		logger:    log.NewEntry(log.New()),
		clientset: client,
		informer:  informer,
		queue:     queue,
		handler:   handlerc,
		name:      obj_type,
	}

	return &controller
}

func Synchronizer() {
	client := GetKubernetesClient()
	for {
		time.Sleep(time.Hour)
		handlers.NamespaceSynchronize(client)
		handlers.StatefulSetSynchronize(client)
		handlers.DeploymentSynchronize(client)
		handlers.ReplicaSetSynchronize(client)
		handlers.PodSynchronize(client)
		handlers.ServiceSynchronize(client)
		handlers.IngressSynchronize(client)
	}
}

// main code path
func main() {

	log.Info("Current namespace: ", os.Getenv("APP_NAMESPACE"))
	podcontroller := CreateController("Pod")
	svccontroller := CreateController("Service")
	nscontroller := CreateController("Namespace")
	depcontroller := CreateController("Deployment")
	rscontroller := CreateController("ReplicaSet")
	ingcontroller := CreateController("Ingress")
	sscontroller := CreateController("StatefulSet")

	// use a channel to synchronize the finalization for a graceful shutdown
	stopCh := make(chan struct{})
	defer close(stopCh)

	// log.SetLevel(log.DebugLevel)

	// start sync task
	go Synchronizer()
	// run the controller loop to process items
	go podcontroller.Run(stopCh)
	go svccontroller.Run(stopCh)
	go nscontroller.Run(stopCh)
	go depcontroller.Run(stopCh)
	go rscontroller.Run(stopCh)
	go ingcontroller.Run(stopCh)
	go sscontroller.Run(stopCh)
	// use a channel to handle OS signals to terminate and gracefully shut
	// down processing
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
