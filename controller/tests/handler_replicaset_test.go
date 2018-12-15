package tests

import (
	"context"
	"testing"
	"time"

	"github.com/intuit/katlas/controller/handlers"
	v1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type ReplicaSetTest struct {
	in  *v1beta2.ReplicaSet
	out map[string]interface{}
}

var replicasettests = []ReplicaSetTest{
	{
		in: &v1beta2.ReplicaSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "test-replicaset",
				Namespace:       "test-namespace",
				ResourceVersion: "1",
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind: "Deployment",
						Name: "test-deployment",
					},
				},
				CreationTimestamp: metav1.Time{},
			},
			Spec: v1beta2.ReplicaSetSpec{
				Replicas: replicaptr,
			},
		},
		out: map[string]interface{}{
			"objtype":         "ReplicaSet",
			"name":            "test-replicaset",
			"creationtime":    metav1.Time{},
			"namespace":       "test-namespace",
			"numreplicas":     replicaptr,
			"podspec":         v1.PodSpec{},
			"owner":           "test-deployment",
			"cluster":         "test2",
			"resourceversion": "1",
			"labels":          map[string]string{},
			"k8sobj":          "K8sObj",
		},
	},
}

func TestReplicaSet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = ctx
	_ = cancel

	// Create the fake client.
	client := fake.NewSimpleClientset()
	// replicasetcontroller := controller.CreateController("ReplicaSet", client)
	// go replicasetcontroller.Run(stopCh)
	replicasethandler := handlers.ReplicaSetHandler{}
	replicasethandler.Init()

	for i, test := range replicasettests {
		_ = i
		testreplicaset := test.in
		a, err := client.AppsV1beta2().ReplicaSets("test-namespace").Create(testreplicaset)
		if err != nil {
			t.Errorf("error injecting replicaset add: %v", err)
		}

		listreplicasets, err := client.AppsV1beta2().ReplicaSets("test-namespace").List(metav1.ListOptions{})

		t.Logf("ReplicaSets: %s\n", listreplicasets.String())

		output, err := replicasethandler.ObjectCreated(a)
		if err != nil {
			t.Errorf("error reading replicaset : %v", err)
		}
		err = replicasethandler.ObjectUpdated(a, a)

		if len(output) != len(test.out) {
			t.Error("Informer did not get the added replicaset diff length")
		}

		for k, v := range output {
			if k == "podspec" || k == "numreplicas" || k == "labels" || k == "creationtime" {
				continue
			}
			if w, ok := test.out[k]; !ok || v != w {
				t.Error("Informer did not get the added replicaset diff values")
				t.Errorf("output[%s]: %s\n", k, output[k])
				t.Errorf("output[%s]: %s\n", k, test.out[k])
			}
		}

		t.Logf("Got replicaset from channel: %s/%s\n", output["replicaset"], test.out["name"])

	}

	handlers.ReplicaSetSynchronize(client)
	t.Log("ReplicaSets synced")

	replicasetinformer := handlers.GetReplicaSetInformer(client)
	if replicasetinformer == nil {
		t.Error("error creating replicaset informer")
	}
}
