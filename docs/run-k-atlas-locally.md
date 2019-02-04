# Run K-Atlas locally

1. Clone the repository.
2. Run Dgraph [https://tour.dgraph.io/intro/2/](https://tour.dgraph.io/intro/2/)
3. Setup Dgraph schema locally \(TODO- We will have script/code for this\)
4. Run the K-Atlas Service
    a) Run from code:
        * cd katlas/service/
        * Get the necessary dependencies

        ```text
        go run server.go
        ```

    b) Run from built docker image:
        * cd katlas/service/ && make all
        * docker build --no-cache -f Dockerfile -t katlas/katlas-service .
        * docker run -d -it -p 8011:8011 -e ENV_NAMESPACE=qal -e SERVER_TYPE=http -e DGRAPH_HOST=$HOST_IP_ADDRESS:9080 --name katlas-service katlas/katlas-service

5. Run the Collector

* cd katlas/controller/
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
$ cd katlas/app/
$ yarn install
$ yarn start
```

This will open the Web viewer in the browser



