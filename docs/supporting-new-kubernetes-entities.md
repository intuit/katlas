# Supporting new Kubernetes Entities

#### Currently Supported Kubernetes Entities

Below are the currently supported Kubernetes Entities that the Collector will collect information about based on the cluster it is deployed in. Support to track additional Kubernetes entities is possible with some code changes as listed below. For contributing your changes,  please refer to the [Contributing guidelines](contributing.md).

* Pod
* Deployment
* Service
* Ingress
* StatefulSet
* ReplicaSet

#### Tracking additional Kubernetes object types

1. Duplicate one of the files under cutlass/controllers/handlers/ and rename it handler\_{NewObjectType}.go.
2. Replace instances of the old object type and api version with the new object type and its appropriate api version
3. Define the object metadata to be extracted in Create{ObjectType}Data\(\)
4. In main.go, add an additional case in the switch statement in the CreateController function to the informer and handler of the new object type
5. In main.go, call CreateController with your object type and run the controller as a goroutine in the main function
6. If necessary, ensure that the controller's service account has a role with the ability to list/watch/get the new object type. Add it as a resource in the cluster role file and apply it against the cluster or bind the service account to a role that has the permissions that you need.

Tests for new handlers can reflect the other existing tests as well.

#### Modifying data extracted from object types

Each handler has a Create{ObjectType}Query function that is responsible for extracting the data from the kubernetes metadata and returning a map to be sent to the rest service. Modifying this map will modify the data sent.

