package tests

import (
	"testing"

	"github.com/intuit/katlas/controller/handlers"
)

type Test struct {
	in  map[string]interface{}
	out int
}

var tests = []Test{
	{
		map[string]interface{}{"clustername": "paas-preprod-west2.cluster.k8s.local",
			"name":            "dev-test-tep-service-usw2-ppd-e2e",
			"objtype":         "Namespace",
			"resourceversion": "9162214"},
		200,
	},

	{
		map[string]interface{}{"availablereplicas": 1,
			"clustername":     "paas-preprod-west2.cluster.k8s.local",
			"creationtime":    "2018-10-31T20:12:18Z",
			"labels":          map[string]interface{}{"app": "tep-ws"},
			"name":            "tep-ws-deployment-just-stash",
			"namespace":       "dev-test-tep-service-usw2-ppd-qal",
			"numreplicas":     1,
			"objtype":         "Deployment",
			"resourceversion": "25441337",
			"strategy":        "RollingUpdate"},
		200,
	},
	{
		map[string]interface{}{"clustername": "paas-preprod-west2.cluster.k8s.local",
			"containers":      []map[string]interface{}{{"name": "operator-nginx", "image": "nginx", "ports": []map[string]interface{}{{"containerPort": 80, "protocol": "TCP"}}, "resources": map[string]interface{}{}, "volumeMounts": []map[string]interface{}{{"name": "default-token-s8lrn", "readOnly": true, "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"}}, "terminationMessagePath": "/dev/termination-log", "terminationMessagePolicy": "File", "imagePullPolicy": "Always"}},
			"ip":              "",
			"labels":          map[string]interface{}{"pod-template-hash": "3003317022", "run": "operator-nginx"},
			"name":            "operator-nginx-744775c466-l74xr",
			"namespace":       "monitoring",
			"nodename":        "ip-10-83-122-52.us-west-2.compute.internal",
			"objtype":         "Pod",
			"owner":           "operator-nginx-744775c466",
			"ownertype":       "ReplicaSet",
			"phase":           "Pending",
			"resourceversion": "28371618",
			"starttime":       nil,
			"volumes":         []map[string]interface{}{{"name": "default-token-s8lrn", "secret": map[string]interface{}{"secretName": "default-token-s8lrn", "defaultMode": 420}}}},
		200,
	},
	{
		map[string]interface{}{"clustername": "paas-preprod-west2.cluster.k8s.local",
			"creationtime": "2018-10-23T13:21:58Z",
			"labels": map[string]interface{}{
				"app":               "k8-sample-v4-service",
				"pod-template-hash": "3185861569",
				"release":           "k8-sample-v4-service",
				"splunk-index":      "k8-sample-v4-service-qal"},
			"name":        "k8-sample-v4-service-75d9db59bf",
			"namespace":   "dev-patterns-k8-sample-v4-service-usw2-ppd-qal",
			"numreplicas": 0,
			"objtype":     "ReplicaSet",
			"owner":       "k8-sample-v4-service",
			"podspec": map[string]interface{}{
				"volumes": []map[string]interface{}{
					{"name": "nginx-secret-volume", "secret": map[string]interface{}{"secretName": "nginxsecret", "defaultMode": 420}},
					{"name": "nginx-conf-volume", "configMap": map[string]interface{}{"name": "nginxconfigmap", "defaultMode": 420}},
					{"name": "logs-volume", "emptyDir": map[string]interface{}{}},
					{"name": "splunk-config", "configMap": map[string]interface{}{"name": "splunk-forwarder-config", "defaultMode": 420}}},
				"containers": []map[string]interface{}{
					{"name": "k8-sample-v4-service",
						"image": "docker.artifactory.a.com/dev/patterns/k8-sample-v4-service/service/k8-sample-v4-service:0bfd4f8",
						"ports": []map[string]interface{}{{"containerPort": 8080, "protocol": "TCP"}},
						"env": []map[string]interface{}{
							{"name": "APPDYNAMICS_AGENT_ACCOUNT_NAME", "value": "ss-dev"},
							{"name": "APPDYNAMICS_CONTROLLER_HOST_NAME", "value": "ss-dev.saas.appdynamics.com"},
							{"name": "APPDYNAMICS_CONTROLLER_PORT", "value": "443"},
							{"name": "APPDYNAMICS_AGENT_APPLICATION_NAME", "value": "k8-push-event-processor-service"},
							{"name": "APPDYNAMICS_AGENT_NODE_NAME", "value": "k8-node-push-event-processor-service"},
							{"name": "APPDYNAMICS_AGENT_TIER_NAME", "value": "app"},
							{"name": "MACHINE_AGENT_PROPERTIES", "value": "-Dappdynamics.sim.enabled=true -Dappdynamics.docker.enabled=false"},
							{"name": "APPDYNAMICS_AGENT_ACCOUNT_ACCESS_KEY", "valueFrom": map[string]interface{}{"secretKeyRef": map[string]interface{}{"name": "appd-secret", "key": "access-key"}}}},
						"resources":    map[string]interface{}{},
						"volumeMounts": []map[string]interface{}{{"name": "logs-volume", "mountPath": "/OILLogs"}},
						"livenessProbe": map[string]interface{}{"httpGet": map[string]interface{}{"path": "/health/local", "port": 8080, "scheme": "HTTP"},
							"initialDelaySeconds": 60,
							"timeoutSeconds":      1,
							"periodSeconds":       60,
							"successThreshold":    1,
							"failureThreshold":    3},
						"readinessProbe": map[string]interface{}{"httpGet": map[string]interface{}{"path": "/health/local", "port": 8080, "scheme": "HTTP"},
							"initialDelaySeconds": 60,
							"timeoutSeconds":      1,
							"periodSeconds":       60,
							"successThreshold":    1,
							"failureThreshold":    3},
						"terminationMessagePath":   "/dev/termination-log",
						"terminationMessagePolicy": "File",
						"imagePullPolicy":          "IfNotPresent"},
					{"name": "nginx",
						"image":                    "nginx:1.13",
						"command":                  []string{"sh", "-c"},
						"args":                     []string{"cp /conf/nginx.conf /etc/nginx/nginx.conf \\u0026\\u0026 nginx -g 'daemon off;'"},
						"ports":                    []map[string]interface{}{{"containerPort": 80, "protocol": "TCP"}, {"containerPort": 443, "protocol": "TCP"}},
						"resources":                map[string]interface{}{},
						"volumeMounts":             []map[string]interface{}{{"name": "nginx-conf-volume", "readOnly": true, "mountPath": "/conf"}, {"name": "nginx-secret-volume", "readOnly": true, "mountPath": "/etc/nginx/ssl"}},
						"terminationMessagePath":   "/dev/termination-log",
						"terminationMessagePolicy": "File",
						"imagePullPolicy":          "IfNotPresent"},
					{"name": "splunkuf",
						"image": "docker.artifactory.a.com/dev/patterns/k8-sample-v4-service/service/splunk-forwarder/splunkforwarder-debian-9:latest",
						"ports": []map[string]interface{}{{"containerPort": 80, "protocol": "TCP"}},
						"env": []map[string]interface{}{{"name": "SPLUNK_START_ARGS", "value": "--accept-license"}, {"name": "SPLUNK_PASSWORD", "value": "123456789"}, {"name": "IDSP_POLICY_ID", "value": "p-hxq68k5gm9v9"},
							{"name": "IDPS_ENDPOINT",
								"value": "PatternAutomation-PRE-PRODUCTION-KJUEA7.pd.idps.a.com"}},
						"resources":                map[string]interface{}{},
						"volumeMounts":             []map[string]interface{}{{"name": "logs-volume", "mountPath": "/applogs"}, {"name": "splunk-config", "mountPath": "/splunk-app/defaults"}},
						"terminationMessagePath":   "/dev/termination-log",
						"terminationMessagePolicy": "File", "imagePullPolicy": "Always"}},
				"restartPolicy":                 "Always",
				"terminationGracePeriodSeconds": 30,
				"dnsPolicy":                     "ClusterFirst",
				"securityContext":               map[string]interface{}{},
				"schedulerName":                 "default-scheduler"},
			"resourceversion": "22447313"},
		200,
	},
}

// func TestSendJSONQuery(t *testing.T) {
// 	for i, test := range tests {
// 		code, data := handlers.SendJSONQuery(test.in, "http://localhost:8011/v1/entity/"+test.in["name"].(string))
// 		_ = data
// 		_ = i
// 		if code != 200 {
// 			//t.Errorf("failed to send jsonquery with namespace data")
// 			t.Error(string(data))
// 		} else {
// 			t.Log(string(data))
// 		}
// 	}
//
// }

func TestGetKubernetesClient(t *testing.T) {
	client := handlers.GetKubernetesClient()
	if client == nil {
		t.Errorf("failed to create client")
	}

}
