package tests

import (
	"context"
	"testing"
	"time"

	"github.com/intuit/katlas/controller/handlers"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type PodTest struct {
	in  *v1.Pod
	out map[string]interface{}
}

var testcluster = "testcluster"

var podtests = []PodTest{
	{
		in: &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "test-pod",
				Namespace:       "test-namespace",
				ResourceVersion: "1",
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind: "ReplicaSet",
						Name: "test-replicaset",
					},
				},
				CreationTimestamp: *timeptr,
			},
			Status: v1.PodStatus{
				Phase:     v1.PodPending,
				StartTime: &metav1.Time{},
				PodIP:     "192.168.1.1",
			},
			Spec: v1.PodSpec{
				NodeName: "node1",
			},
		},
		out: map[string]interface{}{
			"objtype":         "pod",
			"namespace":       "test-namespace",
			"creationtime":    *timeptr,
			"phase":           v1.PodPending,
			"ip":              "192.168.1.1",
			"containers":      []v1.Container{},
			"labels":          map[string]string{},
			"name":            "test-pod",
			"resourceversion": "1",
			"nodename":        "node1",
			"volumes":         []v1.Volume{},
			"owner":           "test-replicaset",
			"ownertype":       "ReplicaSet",
			"cluster":         "test2",
			"k8sobj":          "K8sObj",
		},
	},
}

func TestPod(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = ctx
	_ = cancel

	// Create the fake client.
	client := fake.NewSimpleClientset()
	// podcontroller := controller.CreateController("Pod", client)
	// go podcontroller.Run(stopCh)
	podhandler := handlers.PodHandler{}
	podhandler.Init()

	for i, test := range podtests {
		_ = i
		testpod := test.in
		a, err := client.Core().Pods("test-namespace").Create(testpod)
		if err != nil {
			t.Errorf("error injecting pod add: %v", err)
		}

		listpods, err := client.Core().Pods("test-namespace").List(metav1.ListOptions{})
		t.Logf("Pods: %s\n", listpods.String())

		err = podhandler.ObjectCreated(a)
		if err != nil {
			t.Errorf("error creating pod : %v", err)
		}
		err = podhandler.ObjectUpdated(a, a)
		if err != nil {
			t.Errorf("error updating pod : %v", err)
		}
	}

	handlers.PodSynchronize(client)
	t.Log("Pods synced")

	podinformer := handlers.GetPodInformer(client)
	if podinformer == nil {
		t.Error("error creating pod informer")
	}

}
