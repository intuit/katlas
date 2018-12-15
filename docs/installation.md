# Installation

With an introduction to the core concepts in the previous sections, lets move onto setting up K-Atlas on a Kubernetes cluster and using it.

## Prerequisites

**The below installation will assume a Kubernetes cluster with Minikube is being used.**

1. Install Minikube

{% embed url="https://github.com/kubernetes/minikube" %}

2. Once minikube is up and running, Note that cluster `minikube`  and namespace `default` exists after the above steps. The yaml files used in the Installation use this cluster and namespace.

{% hint style="info" %}
If using another cluster and namespace, the K-Atlas Installation yaml files will need to be modified to use the correct cluster and namespace.
{% endhint %}

```text
$ kubectl config get-clusters
NAME
minikube
 
$ kubectl get namespace
NAME          STATUS   AGE
default       Active   11m
kube-public   Active   11m
kube-system   Active   11m
```

Then move onto installing K-Atlas. 

## Install K-Atlas

The individual components must be installed in the below order-

* Dgraph
* K-Atlas Service
* K-Atlas Kubernetes Collector
* K-Atlas Browser

```text
$ kubectl create -f deploy/dgraph.yaml

$ kubectl get services
NAME            TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                      AGE
dgraph-public   LoadBalancer   10.109.53.229   <pending>     5080:30766/TCP,6080:32699/TCP,8080:32038/TCP,9080:30796/TCP,8000:31572/TCP   1d

$ kubectl get pods
NAME                                 READY   STATUS    RESTARTS   AGE
dgraph-0                             3/3     Running   0          1d
```

Setup Dgraph with the Schema and Metadata \[TODO- This step is temporary, after more code changes, it will be done by the code\]

```text
cd deploy/
go run create_db_schema.go -dbhost=<> -port=<>
go run create_k8smeta.go -dbhost=<> -port=<>

where dbhost = minikube-ip
where port = NodePort corresponding to Port 9080, in the output of 'kubectl get service dgraph-public' .
```

Install K-Atlas Service

```text
$ kubectl create -f deploy/katlas-service.yaml
```

Install K-Atlas Collector

```text
$ kubectl create -f deploy/katlas-collector.yaml
```

Install K-Atlas Browser

{% hint style="info" %}
Modify the below env variable in deploy/katlas-browser.yaml and then run the kubectl command.

KATLAS\_API\_URL

value: http://&lt;minikube-ip&gt;:30415
{% endhint %}

```text
$ kubectl create -f deploy/katlas-browser.yaml
```

Check all pods are up

```text
$ kubectl get pods
NAME                                 READY   STATUS    RESTARTS   AGE
dgraph-0                             3/3     Running   9          1d
katlas-service-748668b795-wt6gk      1/1     Running   3          1d
katlas-controller-8586d84564-njrvr   1/1     Running   3          1d
katlas-browser-5875c79c64-2zhwk      1/1     Running   3          1d
```

Check all Services are running

```text
$ kubectl get services
dgraph-public   LoadBalancer   10.105.204.11    <pending>     5080:30252/TCP,6080:31104/TCP,8080:32395/TCP,9080:30796/TCP,8000:30063/TCP   1d
katlas-service  LoadBalancer   10.96.175.179    <pending>     8011:30415/TCP                                                               1d
katlas-browser  LoadBalancer   10.106.140.250   <pending>     80:30417/TCP                                                                 1d
```

If using minikube, point your browser to the following URL to start using K-Atlas Browser:

```text
http://<minikube-ip>:30417
```

{% hint style="info" %}
Ensure you use the Chrome browser. For issues with Installation, please refer to the FAQ Section.
{% endhint %}



