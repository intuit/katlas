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
			"objtype":         "namespace",
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

		err = namespacehandler.ObjectCreated(a)
		if err != nil {
			t.Errorf("error creating namespace : %v", err)
		}
		err = namespacehandler.ObjectUpdated(a, a)
		if err != nil {
			t.Errorf("error updating namespace : %v", err)
		}
	}

	handlers.NamespaceSynchronize(client)
	t.Log("Namespaces synced")

	namespaceinformer := handlers.GetNamespaceInformer(client)
	if namespaceinformer == nil {
		t.Error("error creating namespace informer")
	}
}
