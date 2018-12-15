set -x  #echo on

kubectl delete -f katlas-browser.yaml
kubectl delete -f katlas-collector.yaml
kubectl delete -f katlas-service.yaml
kubectl delete -f dgraph.yaml
