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
      "name": "preprod-west2.cluster.k8s.local",
      "objtype": "Cluster",
      "resourceid": "preprod-west2.cluster.k8s.local",
      "resourceversion": "0",
      "uid": "0x6f3"
    }],
    "creationtime": "2018-11-30T01:06:10Z",
    "k8sobj": "K8sObj",
    "labels": "null",
    "name": "katlas-browser-deployment",
    "namespace": [{
      "k8sobj": "K8sObj",
      "labels": "{\"owner\":\"iksm\",\"prune-label\":\"katlas-usw2-ppd-qal\",\"name\":\"katlas-usw2-ppd-qal\"}",
      "name": "katlas-usw2-ppd-qal",
      "objtype": "Namespace",
      "resourceid": "preprod-west2.cluster.k8s.local:katlas-usw2-ppd-qal",
      "resourceversion": "29971747",
      "uid": "0x759"
    }],
    "numreplicas": 1,
    "objtype": "Deployment",
    "resourceid": "preprod-west2.cluster.k8s.local:katlas-usw2-ppd-qal:katlas-browser-deployment",
    "resourceversion": "47750611",
    "strategy": "RollingUpdate",
    "uid": "0xafbe",
    "~owner": [{
      "creationtime": "2018-11-30T01:06:10Z",
      "k8sobj": "K8sObj",
      "labels": "{\"app\":\"katlas-browser\",\"pod-template-hash\":\"1588138086\",\"splunk-index\":\"k8s_paas\"}",
      "name": "katlas-browser-deployment-59dd57d4db",
      "numreplicas": 0,
      "objtype": "ReplicaSet",
      "podspec": "{\"containers\":[{\"env\":[{\"name\":\"ENV_NAMESPACE\",\"value\":\"katlas-usw2-ppd-qal\"}],\"image\":\"docker.artifactory.com/k8scmdb/service/katlas-browser:d14bc38\",\"imagePullPolicy\":\"Always\",\"name\":\"katlas-browser\",\"ports\":[{\"containerPort\":80,\"protocol\":\"TCP\"}],\"resources\":{},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\"}],\"dnsPolicy\":\"ClusterFirst\",\"restartPolicy\":\"Always\",\"schedulerName\":\"default-scheduler\",\"securityContext\":{},\"terminationGracePeriodSeconds\":30}",
      "resourceid": "preprod-west2.cluster.k8s.local:katlas-usw2-ppd-qal:katlas-browser-deployment-59dd57d4db",
      "resourceversion": "37599218",
      "uid": "0xafbf"
    }, {
      "creationtime": "2018-11-30T17:59:40Z",
      "k8sobj": "K8sObj",
      "labels": "{\"app\":\"katlas-browser\",\"pod-template-hash\":\"3490095556\",\"splunk-index\":\"k8s_paas\"}",
      "name": "katlas-browser-deployment-78f44f999b",
      "numreplicas": 0,
      "objtype": "ReplicaSet",
      "podspec": "{\"containers\":[{\"env\":[{\"name\":\"ENV_NAMESPACE\",\"value\":\"katlas-usw2-ppd-qal\"}],\"image\":\"docker.artifactory.com/k8scmdb/service/katlas-browser:4fc1f5f\",\"imagePullPolicy\":\"Always\",\"name\":\"katlas-browser\",\"ports\":[{\"containerPort\":80,\"protocol\":\"TCP\"}],\"resources\":{},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\"}],\"dnsPolicy\":\"ClusterFirst\",\"restartPolicy\":\"Always\",\"schedulerName\":\"default-scheduler\",\"securityContext\":{},\"terminationGracePeriodSeconds\":30}",
      "resourceid": "preprod-west2.cluster.k8s.local:katlas-usw2-ppd-qal:katlas-browser-deployment-78f44f999b",
      "resourceversion": "38670868",
      "uid": "0xbb53"
    }, {
      "creationtime": "2018-11-30T17:59:40Z",
      "k8sobj": "K8sObj",
      "labels": "{\"app\":\"katlas-browser\",\"pod-template-hash\":\"3490095556\",\"splunk-index\":\"k8s_paas\"}",
      "name": "katlas-browser-deployment-78f44f999b",
      "numreplicas": 1,
      "objtype": "ReplicaSet",
      "podspec": "{\"containers\":[{\"env\":[{\"name\":\"ENV_NAMESPACE\",\"value\":\"katlas-usw2-ppd-qal\"}],\"image\":\"docker.artifactory.com/k8scmdb/service/katlas-browser:4fc1f5f\",\"imagePullPolicy\":\"Always\",\"name\":\"katlas-browser\",\"ports\":[{\"containerPort\":80,\"protocol\":\"TCP\"}],\"resources\":{},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\"}],\"dnsPolicy\":\"ClusterFirst\",\"restartPolicy\":\"Always\",\"schedulerName\":\"default-scheduler\",\"securityContext\":{},\"terminationGracePeriodSeconds\":30}",
      "resourceid": "preprod-west2.cluster.k8s.local:katlas-usw2-ppd-qal:katlas-browser-deployment-78f44f999b",
      "resourceversion": "37698762",
      "uid": "0xbb55"
    }]
  }],
  "obj5": []
};