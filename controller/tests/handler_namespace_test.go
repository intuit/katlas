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

type NamespaceTest struct {
	in  *v1.Namespace
	out map[string]interface{}
}

var namespacetests = []NamespaceTest{
	{
		in: &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "test-namespace",
				ResourceVersion:   "1",
				Labels:            map[string]string{"label1": "value1"},
				CreationTimestamp: *timeptr,
			},
		},
		out: map[string]interface{}{
			"objtype":         "Namespace",
			"name":            "test-namespace",
			"resourceversion": "1",
			"cluster":         "test2",
			"k8sobj":          "K8sObj",
			"labels":          map[string]string{"label1": "value1"},
			"creationtime":    *timeptr,
		},
	},
}

func TestNamespace(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = ctx
	_ = cancel

	// Create the fake client.
	client := fake.NewSimpleClientset()
	// namespacecontroller := controller.CreateController("Namespace", client)
	// go namespacecontroller.Run(stopCh)
	namespacehandler := handlers.NamespaceHandler{}
	namespacehandler.Init()

	for i, test := range namespacetests {
		_ = i
		testnamespace := test.in
		a, err := client.Core().Namespaces().Create(testnamespace)
		if err != nil {
			t.Errorf("error injecting namespace add: %v", err)
		}

		listnamespaces, err := client.Core().Namespaces().List(metav1.ListOptions{})

		t.Logf("Namespaces: %s\n", listnamespaces.String())

		output, err := namespacehandler.ObjectCreated(a)
		if err != nil {
			t.Errorf("error reading namespace : %v", err)
		}
		err = namespacehandler.ObjectUpdated(a, a)

		if len(output) != len(test.out) {
			t.Error("Informer did not get the added namespace diff length")
		}

		for k, v := range output {
			if k == "volumes" || k == "containers" || k == "labels" || k == "creationtime" {
				continue
			}
			if w, ok := test.out[k]; !ok || v != w {
				t.Error("Informer did not get the added namespace diff values")
				t.Errorf("output[%s]: %s\n", k, output[k])
				t.Errorf("output[%s]: %s\n", k, test.out[k])
			}
		}

		t.Logf("Got namespace from channel: %s/%s\n", output["namespace"], test.out["name"])

	}

	handlers.NamespaceSynchronize(client)
	t.Log("Namespaces synced")

	namespaceinformer := handlers.GetNamespaceInformer(client)
	if namespaceinformer == nil {
		t.Error("error creating namespace informer")
	}
}
