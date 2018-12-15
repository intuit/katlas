# Run K-Atlas locally

1. Clone the repository.
2. Run Dgraph [https://tour.dgraph.io/intro/2/](https://tour.dgraph.io/intro/2/)
3. Setup Dgraph schema locally \(TODO- We will have script/code for this\)
4. Run the K-Atlas Service

* cd cutlass/rest-service/
* Get the necessary dependencies

```text
go run server.go
```

    5. Run the Collector

* cd cutlass/controller/
* Get the necessary dependencies
* Ensure that your kubeconfig for the cluster that you want to monitor is at $HOME/.kube/config
* Set necessary environment variables

```text
export CLUSTER_NAME={cluster_name}
export TARGET_URL={k-atlas-api_url}
```

Run the below command to start the collector.

```text
go run *.go
```

    6. Run the Web Viewer

```text
$ cd cutlass/app/
$ yarn install
$ yarn start
```

This will open the Web viewer in the browser



