package tests

import (
	"context"
	"testing"
	"time"

	"github.com/intuit/katlas/controller/handlers"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type StatefulsetTest struct {
	in  *appsv1.StatefulSet
	out map[string]interface{}
}

var replicas = int32(2)

var statefulsettests = []StatefulsetTest{
	{
		in: &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "test-statefulset",
				Namespace:         "test-namespace",
				CreationTimestamp: metav1.Time{},
				ResourceVersion:   "1",
			},
			Spec: appsv1.StatefulSetSpec{
				Replicas: &replicas,
			},
		},
		out: map[string]interface{}{
			"objtype":         "StatefulSet",
			"name":            "test-statefulset",
			"namespace":       "test-namespace",
			"creationtime":    metav1.Time{},
			"numreplicas":     &replicas,
			"cluster":         "test2",
			"resourceversion": "1",
			"labels":          map[string]string{},
			"k8sobj":          "K8sObj",
		},
	},
}

func TestStatefulset(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = ctx
	_ = cancel

	// Create the fake client.
	client := fake.NewSimpleClientset()

	statefulsethandler := handlers.StatefulSetHandler{}
	statefulsethandler.Init()

	for i, test := range statefulsettests {
		_ = i
		teststatefulset := test.in
		a, err := client.AppsV1().StatefulSets("test-namespace").Create(teststatefulset)
		if err != nil {
			t.Errorf("error injecting statefulset add: %v", err)
		}

		liststatefulsets, err := client.AppsV1().StatefulSets("test-namespace").List(metav1.ListOptions{})
		t.Logf("Statefulsets: %s\n", liststatefulsets.String())

		output, err := statefulsethandler.ObjectCreated(a)
		if err != nil {
			t.Errorf("error reading statefulset : %v", err)
		}
		err = statefulsethandler.ObjectUpdated(a, a)

		if len(output) != len(test.out) {
			t.Error("Informer did not get the added statefulset diff length")
		}

		for k, v := range output {
			if k == "selector" || k == "numreplicas" || k == "labels" || k == "creationtime" {
				continue
			}
			if w, ok := test.out[k]; !ok || v != w {
				t.Error("Informer did not get the added statefulset diff values")
				t.Errorf("output[%s]: %s\n", k, output[k])
				t.Errorf("test.out[%s]: %s\n", k, test.out[k])
			}
		}

		t.Logf("Got statefulset from channel: %s/%s\n", output["namespace"], test.out["name"])

	}

	handlers.StatefulSetSynchronize(client)
	t.Log("StatefulSets synced")

	statefulsetinformer := handlers.GetStatefulSetInformer(client)
	if statefulsetinformer == nil {
		t.Error("error creating statefulset informer")
	}

}
