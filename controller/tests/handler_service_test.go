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

type ServiceTest struct {
	in  *v1.Service
	out map[string]interface{}
}

var servicetests = []ServiceTest{
	{
		in: &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "test-service",
				Namespace:         "test-namespace",
				ResourceVersion:   "1",
				CreationTimestamp: *timeptr,
			},
			Spec: v1.ServiceSpec{
				Type:      v1.ServiceTypeNodePort,
				ClusterIP: "192.168.1.1",
			},
		},
		out: map[string]interface{}{
			"objtype":         "Service",
			"name":            "test-service",
			"namespace":       "test-namespace",
			"selector":        map[string]string{},
			"labels":          map[string]string{},
			"clusterip":       "192.168.1.1",
			"servicetype":     v1.ServiceTypeNodePort,
			"ports":           []v1.ServicePort{},
			"resourceversion": "1",
			"cluster":         "test2",
			"k8sobj":          "K8sObj",
			"creationtime":    *timeptr,
		},
	},
}

func TestService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = ctx
	_ = cancel

	// Create the fake client.
	client := fake.NewSimpleClientset()

	servicehandler := handlers.ServiceHandler{}
	servicehandler.Init()

	for i, test := range servicetests {
		_ = i
		testservice := test.in
		a, err := client.Core().Services("test-namespace").Create(testservice)
		if err != nil {
			t.Errorf("error injecting service add: %v", err)
		}

		listservices, err := client.Core().Services("test-namespace").List(metav1.ListOptions{})
		t.Logf("Services: %s\n", listservices.String())

		output, err := servicehandler.ObjectCreated(a)
		if err != nil {
			t.Errorf("error reading service : %v", err)
		}
		err = servicehandler.ObjectUpdated(a, a)

		if len(output) != len(test.out) {
			t.Error("Informer did not get the added service diff length")
		}

		for k, v := range output {
			if k == "selector" || k == "ports" || k == "labels" {
				continue
			}
			if w, ok := test.out[k]; !ok || v != w {
				t.Error("Informer did not get the added service diff values")
				t.Errorf("output[%s]: %s\n", k, output[k])
				t.Errorf("test.out[%s]: %s\n", k, test.out[k])
			}
		}

		t.Logf("Got service from channel: %s/%s\n", output["namespace"], test.out["name"])

	}

	handlers.ServiceSynchronize(client)
	t.Log("Services synced")

	serviceinformer := handlers.GetServiceInformer(client)
	if serviceinformer == nil {
		t.Error("error creating service informer")
	}

}
