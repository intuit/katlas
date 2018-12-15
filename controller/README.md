# Cutlass Controller

## Purpose

A controller responsible for collecting and sending information about certain Kinds of Kubernetes objects (Deployments, Ingresses, Namespaces, Pods, ReplicaSets, Services, StatefulSets) to the rest service.

## Function
This controller spawns a thread for each object type being tracked. First, all objects of a given type are listed and certain fields extracted from the Kubernetes object metadata and sent to the rest service as json. Then the thread watches for events of that type and sends the changes when objects are created, modified, or deleted.


## Installation

### To run on local machine

```
1. go get the necessary dependencies
$ go get -u -v k8s.io/client-go/...
$ go get -u -v k8s.io/api/...
$ go get -u -v github.com/dgraph-io/dgo
$ go get -u -v github.com/Sirupsen/logrus

2. in the controller directory, run
$ go run *.go
```
### To run on cluster
```
use a provided deployment file in helm/templates
in controller-deployment-{env}-{region}.yaml
set TARGET_URL to the address of the api
set CLUSTER_NAME to the name of the cluster where the controller is deployed
$ kubectl apply -f helm/templates/controller-deployment-ppd-w.yaml
$ kubectl get pods // to find the pod name of the controller
$ kubectl logs -f cmk-controller-xxxxxxxxx-xxxxx
```

### Running Tests
```
1. set necessary environment variables
$ export CLUSTER_NAME=test2
$ export TARGET_URL=http://localhost:8011/
2. start the test server
$ go run testserver/test-server.go
3. run the tests
$ go test -v -coverpkg ../handlers -coverprofile cover.out
4. convert to html
$ go tool cover -html=cover.out -o cover.html
5. open cover.html in a browser to see coverage percentages and lines
```

The test server is a simple net/http server with endpoints for entity, query, sync and health, to mimic the operation of the actual rest service. The entity endpoint will respond with 200 when well formed data is received and 500 if there was some error handling the json in the request. The Sync endpoint will return a list of objects, some will exist in the client's fake cluster, others will not, and the handler test will send the appropriate delete requests to match its fake cluster.

## Contributing

Visit the [Contribution Documentation]

### Tracking additional object types
1. Duplicate one of the files under handlers/ and rename it handler_{NewObjectType}.go.
2. Replace instances of the old object type and api version with the new object type and appropriate api version
3. Define the object metadata to be extracted in Create{ObjectType}Data()
4. in main.go, add an additional case in CreateController to return your informer and handler
5. in main.go, call CreateController with your object type and run the controller as a goroutine

Tests for new handlers can reflect the other existing tests as well.

### Modifying data extracted from object types
Each handler has a Create{ObjectType}Query function that is responsible for extracting the data from the kubernetes metadata and returning a map to be sent to the rest service. Modifying this map will modify the data sent.
