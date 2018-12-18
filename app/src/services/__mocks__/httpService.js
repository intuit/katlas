export class HttpService  {
  static get({url, params, arrayParams}) {
    return new Promise((resolve, reject) => resolve(resp));
  }
}

const resp = {
  "obj0": [],
  "<snip>": [],
  "obj4": [{
    "availablereplicas": 1,
    "cluster": [{
      "k8sobj": "K8sObj",
      "name": "paas-preprod-west2.cluster.k8s.local",
      "objtype": "Cluster",
      "resourceid": "paas-preprod-west2.cluster.k8s.local",
      "resourceversion": "0",
      "uid": "0x6f3"
    }],
    "creationtime": "2018-11-30T01:06:10Z",
    "k8sobj": "K8sObj",
    "labels": "null",
    "name": "cutlass-ui-deployment",
    "namespace": [{
      "k8sobj": "K8sObj",
      "labels": "{\"iks.intuit.com/owner\":\"iksm\",\"iks.intuit.com/prune-label\":\"dev-devx-cmdb-api-usw2-ppd-qal\",\"name\":\"dev-devx-cmdb-api-usw2-ppd-qal\"}",
      "name": "dev-devx-cmdb-api-usw2-ppd-qal",
      "objtype": "Namespace",
      "resourceid": "paas-preprod-west2.cluster.k8s.local:dev-devx-cmdb-api-usw2-ppd-qal",
      "resourceversion": "29971747",
      "uid": "0x759"
    }],
    "numreplicas": 1,
    "objtype": "Deployment",
    "resourceid": "paas-preprod-west2.cluster.k8s.local:dev-devx-cmdb-api-usw2-ppd-qal:cutlass-ui-deployment",
    "resourceversion": "47750611",
    "strategy": "RollingUpdate",
    "uid": "0xafbe",
    "~owner": [{
      "creationtime": "2018-11-30T01:06:10Z",
      "k8sobj": "K8sObj",
      "labels": "{\"app\":\"cutlass-ui\",\"pod-template-hash\":\"1588138086\",\"splunk-index\":\"k8s_paas\"}",
      "name": "cutlass-ui-deployment-59dd57d4db",
      "numreplicas": 0,
      "objtype": "ReplicaSet",
      "podspec": "{\"containers\":[{\"env\":[{\"name\":\"ENV_NAMESPACE\",\"value\":\"dev-devx-cmdb-api-usw2-ppd-qal\"}],\"image\":\"docker.artifactory.a.intuit.com/dev/devx/k8scmdb/service/cutlass-ui:d14bc38\",\"imagePullPolicy\":\"Always\",\"name\":\"cutlass-ui\",\"ports\":[{\"containerPort\":80,\"protocol\":\"TCP\"}],\"resources\":{},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\"}],\"dnsPolicy\":\"ClusterFirst\",\"restartPolicy\":\"Always\",\"schedulerName\":\"default-scheduler\",\"securityContext\":{},\"terminationGracePeriodSeconds\":30}",
      "resourceid": "paas-preprod-west2.cluster.k8s.local:dev-devx-cmdb-api-usw2-ppd-qal:cutlass-ui-deployment-59dd57d4db",
      "resourceversion": "37599218",
      "uid": "0xafbf"
    }, {
      "creationtime": "2018-11-30T17:59:40Z",
      "k8sobj": "K8sObj",
      "labels": "{\"app\":\"cutlass-ui\",\"pod-template-hash\":\"3490095556\",\"splunk-index\":\"k8s_paas\"}",
      "name": "cutlass-ui-deployment-78f44f999b",
      "numreplicas": 0,
      "objtype": "ReplicaSet",
      "podspec": "{\"containers\":[{\"env\":[{\"name\":\"ENV_NAMESPACE\",\"value\":\"dev-devx-cmdb-api-usw2-ppd-qal\"}],\"image\":\"docker.artifactory.a.intuit.com/dev/devx/k8scmdb/service/cutlass-ui:4fc1f5f\",\"imagePullPolicy\":\"Always\",\"name\":\"cutlass-ui\",\"ports\":[{\"containerPort\":80,\"protocol\":\"TCP\"}],\"resources\":{},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\"}],\"dnsPolicy\":\"ClusterFirst\",\"restartPolicy\":\"Always\",\"schedulerName\":\"default-scheduler\",\"securityContext\":{},\"terminationGracePeriodSeconds\":30}",
      "resourceid": "paas-preprod-west2.cluster.k8s.local:dev-devx-cmdb-api-usw2-ppd-qal:cutlass-ui-deployment-78f44f999b",
      "resourceversion": "38670868",
      "uid": "0xbb53"
    }, {
      "creationtime": "2018-11-30T17:59:40Z",
      "k8sobj": "K8sObj",
      "labels": "{\"app\":\"cutlass-ui\",\"pod-template-hash\":\"3490095556\",\"splunk-index\":\"k8s_paas\"}",
      "name": "cutlass-ui-deployment-78f44f999b",
      "numreplicas": 1,
      "objtype": "ReplicaSet",
      "podspec": "{\"containers\":[{\"env\":[{\"name\":\"ENV_NAMESPACE\",\"value\":\"dev-devx-cmdb-api-usw2-ppd-qal\"}],\"image\":\"docker.artifactory.a.intuit.com/dev/devx/k8scmdb/service/cutlass-ui:4fc1f5f\",\"imagePullPolicy\":\"Always\",\"name\":\"cutlass-ui\",\"ports\":[{\"containerPort\":80,\"protocol\":\"TCP\"}],\"resources\":{},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\"}],\"dnsPolicy\":\"ClusterFirst\",\"restartPolicy\":\"Always\",\"schedulerName\":\"default-scheduler\",\"securityContext\":{},\"terminationGracePeriodSeconds\":30}",
      "resourceid": "paas-preprod-west2.cluster.k8s.local:dev-devx-cmdb-api-usw2-ppd-qal:cutlass-ui-deployment-78f44f999b",
      "resourceversion": "37698762",
      "uid": "0xbb55"
    }]
  }],
  "obj5": []
};