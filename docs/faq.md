# FAQ

### How do I get the Port to connect to Dgraph on the minikube cluster to setup the Schema and Metadata?

  
Run the below command and check the NodePort corresponding to port 9080. Here port 30796 will be used.

$ kubectl get services

NAME                  TYPE           CLUSTER-IP       EXTERNAL-IP   PORT\(S\)                         AGE

dgraph-public   LoadBalancer   10.109.53.229   &lt;pending&gt;     5080:30766/TCP,6080:32699/TCP,8080:32038/TCP,9080:30796/TCP,8000:31572/TCP   5h

