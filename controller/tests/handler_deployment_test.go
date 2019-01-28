package tests

import (
	"context"
	"testing"
	"time"

	"github.com/intuit/katlas/controller/handlers"
	"k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type DeploymentTest struct {
	in  *v1beta2.Deployment
	out map[string]interface{}
}

var timeptr = &metav1.Time{Time: time.Now()}

var availablereplicas = int32(2)
var replicaptr = &availablereplicas

var deploymenttests = []DeploymentTest{
	{
		in: &v1beta2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "test-deployment",
				Namespace:         "test-namespace",
				ResourceVersion:   "1",
				CreationTimestamp: *timeptr,
				Labels:            map[string]string{"label1": "value1"},
			},
			Spec: v1beta2.DeploymentSpec{
				Strategy: v1beta2.DeploymentStrategy{
					Type: v1beta2.RollingUpdateDeploymentStrategyType,
				},
				Replicas: replicaptr,
			},
			Status: v1beta2.DeploymentStatus{
				AvailableReplicas: availablereplicas,
			},
		},
		out: map[string]interface{}{
			"objtype":           "deployment",
			"cluster":           "test2",
			"name":              "test-deployment",
			"creationtime":      *timeptr,
			"namespace":         "test-namespace",
			"numreplicas":       replicaptr,
			"availablereplicas": availablereplicas,
			"strategy":          v1beta2.RollingUpdateDeploymentStrategyType,
			"resourceversion":   "1",
			"labels":            map[string]string{"label1": "value1"},
			"k8sobj":            "K8sObj",
		},
	},
}

func TestDeployment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = ctx
	_ = cancel

	// Create the fake client.
	client := fake.NewSimpleClientset()
	stopCh := make(chan struct{})
	defer close(stopCh)
	// deploymentcontroller := controller.CreateController("Deployment", client)
	// go deploymentcontroller.Run(stopCh)

	deploymenthandler := handlers.DeploymentHandler{}
	deploymenthandler.Init()

	for i, test := range deploymenttests {
		_ = i
		testdeployment := test.in
		a, err := client.AppsV1beta2().Deployments("test-namespace").Create(testdeployment)
		if err != nil {
			t.Errorf("error injecting deployment add: %v", err)
		}
		listdeployments, err := client.AppsV1beta2().Deployments("test-namespace").List(metav1.ListOptions{})

		t.Logf("Deployments: %s\n", listdeployments.String())

		err = deploymenthandler.ObjectCreated(a)
		if err != nil {
			t.Errorf("error creating deployment : %v", err)
		}

		err = deploymenthandler.ObjectUpdated(a, a)
		if err != nil {
			t.Errorf("error updating deployment : %v", err)
		}
	}

	handlers.DeploymentSynchronize(client)

	t.Log("Deployments synced")

	deploymentinformer := handlers.GetDeploymentInformer(client)
	if deploymentinformer == nil {
		t.Error("error creating deployment informer")
	}

}
