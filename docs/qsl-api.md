# QSL API

## Purpose

Provide an easy to use language for users to query objects from dgraph without having to use graphQL

## Format
  `objecttype[@filtername="value"]{@field1,@field2}`
  `objecttype1[@filtername="value"]{*}.objecttype2[@fieldname="value",@fieldname2="value2"]{*}`

### Details
  * objecttype - the kubernetes Kind of the object
    * needs to match the case as Kubernetes would
    * e.g. ReplicaSet will work, Replicaset/rEplicaset/etc. will not
    * will automatically capitalize first letter, so pod and Pod will both work
  * filtername,value - filter to get objects of objecttype where the fieldname = value
    * value must be enclosed in quotes and can contain alphanumeric characters, ., -, and _
    * commas can be used as boolean equivalent to AND
    * | can be used as the boolean equivalent to OR
    * AND takes precedence over OR, e.g. a,b,c|d,e === (a&b&c) | (d&e)
    * other comparators (<,>,<=,>=) can be used for data types that support comparison
  * field - the fields of the object that we want to return
    * the "\*" will return all fields
    * n "\*" will get ll fields and relations n edges away from the node
  * objecttype1[...]{\*}.objecttype2[...]{\*} - the . denotes a relationship objecttype1->objecttype2
    this will get all fields from objecttype1 and all objecttype2's related to the results of the first block
    with all their fields
  * objecttype, filtername/value and the fields must be nonempty

## Examples
  ```
  pod[@name="pod1"]{@phase}
    return the value of the phase field for all pods named pod1
  ```

  ```
  cluster[@name="cluster1"]{*}.pod[@name="pod1",@ip="1.1.1.1"]{@phase,@image}
    find all clusters named cluster 1
    then find all pods named pod1 with ip 1.1.1.1 related to those clusters
  ```

  ```
  pod[@name="pod1"]{**}
    return the fields of pod1 and the fields for all direct relationships
  ```

## QSL Queries and Their DGraph Equivalents
  ```
  qsl: cluster[@name="paas-preprod-west2.cluster.k8s.local"]{@name}

  dgraph: query Me($objtype: string, $name: string){
	    objects(func: eq(objtype, Cluster)) @filter(eq(name, paas-preprod-west2.cluster.k8s.local)){
				name
			}
	}
  ```
  ```
  qsl: cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}.namespace[@name="opa"|@name="default"]{*}

  dgraph:   query objects($objtype: string, $name: string){
objects(func: eq(objtype, Cluster)) @filter(( eq(name,"paas-preprod-west2.cluster.k8s.local") )){
	creationtime
	k8sobj
	objtype
	name
	resourceid
	resourceversion
	~cluster @filter(eq(objtype, Namespace) and( eq(name,"opa") or eq(name,"default") )){
		resourceversion
		creationtime
		k8sobj
		labels
		name
		resourceid
		objtype
	}
}
}
  ```

  ```
  qsl: cluster[@name="paas-preprod-west2.cluster.k8s.local"]{@name,@creationtime}.deployment[@name="tiller-deploy"]{@name,@strategy}.ReplicaSet[@name="tiller-deploy-8c8c79584"]{*}

  dgraph:  query objects($objtype: string, $name: string){
objects(func: eq(objtype, Cluster)) @filter( eq(name,"paas-preprod-west2.cluster.k8s.local") ){
	creationtime
	name
	~cluster @filter(eq(objtype, Deployment) and eq(name,"tiller-deploy") ){
		name
		strategy
		~owner @filter(eq(objtype, ReplicaSet) and eq(name,"tiller-deploy-8c8c79584") ){
			labels
			resourceid
			k8sobj
			objtype
			name
			numreplicas
			podspec
			creationtime
			resourceversion
		}
	}
}
}
  ```
  ```
  qsl:  cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}.pod[@name="calico-node-gcj7s"]{*}

  dgraph: query objects($objtype: string, $name: string){
    objects(func: eq(objtype, Cluster)) @filter( eq(name,paas-preprod-west2.cluster.k8s.local)){
      creationtime
      k8sobj
      objtype
      name
      resourceid
      resourceversion
      ~cluster @filter(eq(objtype, Pod) and eq(name,calico-node-gcj7s)){
        resourceid
        k8sobj
        objtype
        name
        labels
        phase
        starttime
        volumes
        ip
        containers
        ownertype
        creationtime
        resourceversion
      }
    }
  }
  ```

  ```
  qsl: localhost:8011/v1/qsl?qslstring=namespace[@name="default"]{**}
  dgraph:  query objects($objtype: string, $name: string){
      objects(func: eq(objtype, Namespace)) @filter( ( eq(name,"default") )){
        expand(_all_){
          expand(_all_){
          }
        }
      }
    }
  ```

## API Responses
### Success

```
input:
namespace[@name="default"]{*}

response:
{
    "objects": [
        {
            "k8sobj": "K8sObj",
            "labels": "null",
            "name": "default",
            "objtype": "Namespace",
            "resourceid": "paas-preprod-west2.cluster.k8s.local:default",
            "resourceversion": "9",
            "uid": "0x714"
        }
    ]
}
```

```
input:
cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}.namespace[@name="opa"]{*}

response:
200 OK
{
    "objects": [
        {
            "k8sobj": "K8sObj",
            "name": "paas-preprod-west2.cluster.k8s.local",
            "objtype": "Cluster",
            "resourceid": "paas-preprod-west2.cluster.k8s.local",
            "resourceversion": "0",
            "~cluster": [
                {
                    "k8sobj": "K8sObj",
                    "labels": "null",
                    "name": "opa",
                    "objtype": "Namespace",
                    "resourceid": "paas-preprod-west2.cluster.k8s.local:opa",
                    "resourceversion": "772"
                }
            ]
        }
    ]
}
```

```
input:
cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}..deployment[@name="dgraph-ratel"]{*}.replicaSet[@name="dgraph-ratel-588856bd5b",@numreplicas>="1"]{*}

response:
200 OK
{
    "objects": [
        {
            "k8sobj": "K8sObj",
            "name": "paas-preprod-west2.cluster.k8s.local",
            "objtype": "Cluster",
            "resourceid": "paas-preprod-west2.cluster.k8s.local",
            "resourceversion": "0",
            "~cluster": [
                {
                    "availablereplicas": 1,
                    "creationtime": "2018-11-20T19:47:56Z",
                    "k8sobj": "K8sObj",
                    "labels": "{\"app\":\"dgraph-ratel\"}",
                    "name": "dgraph-ratel",
                    "numreplicas": 1,
                    "objtype": "Deployment",
                    "resourceid": "paas-preprod-west2.cluster.k8s.local:dev-devx-cmdb-api-usw2-ppd-qal:dgraph-ratel",
                    "resourceversion": "44171156",
                    "strategy": "RollingUpdate"
                },
                {
                    "availablereplicas": 1,
                    "creationtime": "2018-11-27T10:08:49Z",
                    "k8sobj": "K8sObj",
                    "labels": "{\"app\":\"dgraph-ratel\"}",
                    "name": "dgraph-ratel",
                    "numreplicas": 1,
                    "objtype": "Deployment",
                    "resourceid": "paas-preprod-west2.cluster.k8s.local:dev-devx-cmdb-api-usw2-ppd-prf:dgraph-ratel",
                    "resourceversion": "45458408",
                    "strategy": "RollingUpdate",
                    "~owner": [
                        {
                            "creationtime": "2018-11-27T10:08:49Z",
                            "k8sobj": "K8sObj",
                            "labels": "{\"app\":\"dgraph-ratel\",\"pod-template-hash\":\"1444126816\"}",
                            "name": "dgraph-ratel-588856bd5b",
                            "numreplicas": 1,
                            "objtype": "ReplicaSet",
                            "podspec": "{\"containers\":[{\"command\":[\"dgraph-ratel\"],\"image\":\"dgraph/dgraph:v1.0.9\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"ratel\",\"ports\":[{\"containerPort\":8000,\"protocol\":\"TCP\"}],\"resources\":{},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\"}],\"dnsPolicy\":\"ClusterFirst\",\"restartPolicy\":\"Always\",\"schedulerName\":\"default-scheduler\",\"securityContext\":{},\"terminationGracePeriodSeconds\":30}",
                            "resourceid": "paas-preprod-west2.cluster.k8s.local:dev-devx-cmdb-api-usw2-ppd-prf:dgraph-ratel-588856bd5b",
                            "resourceversion": "45458406"
                        }
                    ]
                },
                {
                    "availablereplicas": 1,
                    "creationtime": "2018-12-14T21:20:26Z",
                    "k8sobj": "K8sObj",
                    "labels": "{\"app\":\"dgraph-ratel\"}",
                    "name": "dgraph-ratel",
                    "numreplicas": 1,
                    "objtype": "Deployment",
                    "resourceid": "paas-preprod-west2.cluster.k8s.local:dev-devx-cmdb-api-usw2-ppd-e2e:dgraph-ratel",
                    "resourceversion": "47817014",
                    "strategy": "RollingUpdate"
                }
            ]
        }
    ]
}
```

### Failure
#### Malformed Input
```
input:
cluster[@name="paas-preprod-west2.cluster.k8s.local"]n{*}.namespace[@name="opa"]{*}

response:
400 Bad Request
Malformed Query: cluster[@name="paas-preprod-west2.cluster.k8s.local"]n{*}.namespace[@name="opa"]{*}
```

#### Error Connecting to Dgraph
```
input:
cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}.notrealobject[@name="something"]{*}

response:
500 Internal Server Error
Failed to connect to dgraph to get metadata
```

#### Invalid Relation
```
input:
cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}.notrealobject[@name="something"]{*}

response:
400 Bad Request
no relation found between Nnamespace and Cluster
```
