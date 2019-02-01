package tests

import (
	"context"
	"testing"
	"time"

	"github.com/intuit/katlas/controller/handlers"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type IngressTest struct {
	in  *ext_v1beta1.Ingress
	out map[string]interface{}
}

var ingresstests = []IngressTest{
	{
		in: &ext_v1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "test-ingress",
				Namespace:         "test-namespace",
				ResourceVersion:   "1",
				CreationTimestamp: *timeptr,
			},
			Spec: ext_v1beta1.IngressSpec{},
		},
		out: map[string]interface{}{
			"objtype":         "ingress",
			"cluster":         "test2",
			"name":            "test-ingress",
			"namespace":       "test-namespace",
			"defaultbackend":  nil, //*ext_v1beta1.IngressBackend{},
			"tls":             []ext_v1beta1.IngressTLS{},
			"rules":           []ext_v1beta1.IngressRule{},
			"resourceversion": "1",
			"labels":          map[string]string{},
			"k8sobj":          "K8sObj",
			"creationtime":    *timeptr,
		},
	},
}

func TestIngress(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = ctx
	_ = cancel

	// Create the fake client.
	client := fake.NewSimpleClientset()
	// ingresscontroller := controller.CreateController("Ingress", client)
	// go ingresscontroller.Run(stopCh)
	ingresshandler := handlers.IngressHandler{}
	ingresshandler.Init()

	for i, test := range ingresstests {
		_ = i
		testingress := test.in
		a, err := client.ExtensionsV1beta1().Ingresses("test-namespace").Create(testingress)
		if err != nil {
			t.Errorf("error injecting ingress add: %v", err)
		}

		listingresses, err := client.ExtensionsV1beta1().Ingresses("test-namespace").List(metav1.ListOptions{})

		t.Logf("Ingresses: %s\n", listingresses.String())

		err = ingresshandler.ObjectCreated(a)
		if err != nil {
			t.Errorf("error creating ingress : %v", err)
		}
		err = ingresshandler.ObjectUpdated(a, a)
		if err != nil {
			t.Errorf("error updating ingress : %v", err)
		}
	}

	handlers.IngressSynchronize(client)

	t.Log("Ingresses synced")

	ingressinformer := handlers.GetIngressInformer(client)
	if ingressinformer == nil {
		t.Error("error creating ingress informer")
	}

}
