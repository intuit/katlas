# QSL API

## Purpose

Provide an easy to use language for users to query objects from dgraph without having to use graphQL

## Format
  `objecttype1[@fieldname="value" && @count(relationship)>1 $$first=1,offset=0]{@filedname}.objecttype2[@fieldname="value"]{*}`

### Details
  * objecttype - the kubernetes Kind of the object
    * case-insensitive
    * e.g. ReplicaSet will work, Replicaset/rEplicaset/REPLICASET/etc. will work too
  * `filtername=value` - filter to get objects of objecttype where the filtername = value
    * value must be enclosed in quotes if string type and can contain alphanumeric characters, ., -, and _
    * && can be used as boolean equivalent to AND
    * || can be used as the boolean equivalent to OR
    * AND takes precedence over OR
      * e.g. `a&&b&&c||d&&e` === (a&&b&&c) || (d&&e)
    * other comparators (<,>,<=,>=) can be used for data types that support comparison
      * e.g. `ReplicaSet[@numreplicas>=1]{*}`
  * field - the fields of the object that we want to return
    * each field must begin with an @ followed by an alphanumeric field name, or be a string of \*
     * the list can either only contain comma separated @-prefixed field names or \* strings, not both
    * the "\*" will return all fields of the object in dgraph
    * n "\*" will get all fields and relations n edges away from the node
      * if multiple \* are present any relations after that block are ignored
      * e.g. cluster[...]{\*\*}.namespace[...]{@name} is the same as cluster[...]{\*\*}
  * `objecttype1[...]{\*}.objecttype2[...]{\*}` - the . denotes a relationship objecttype1->objecttype2
    this will get all fields from objecttype1 and all objecttype2's related to the results of the first block
    with all their fields
  * objecttype must be specified
    * filters can be empty and will default to returning all objects of its type
    * fields can also be empty and will default to showing nothing for that object type
  * optional pagination
    * in filters, use `objecttype[@filtername="value"$$first=1,offset=1]{@field1,@field2}`
    * $$first=n will return the first n objects by uid
    * $$offset=m will return the objects in uid order starting from m
    * combine to get $$first=x,offset=y to get the first x objects starting from index y
  * count() supported in filter
    * the function take relationship objtype as parameter and all comparators (=,<,>,<=,>=) can be used

## Examples
  ```
  pod[@name="pod1"]{@phase}
    return the value of the phase field for all pods named pod1
  ```

  ```
  cluster[@name="cluster1"]{*}.pod[@name="pod1" && @ip="1.1.1.1"]{@phase,@image}
    find all clusters named cluster 1
    then find all pods named pod1 with ip 1.1.1.1 related to those clusters
  ```

  ```
  pod[@name="pod1"]{**}
    return all of the fields of pod1 and the fields for all direct relationships
  ```

  ```
  deployment[@name="helm"$$first=10]{*}
    return all of the fields of the first 10 pods named "helm" by uid
  ```

  ```
  ReplicaSet.deployment{*}
    return all of the fields of deployments that have a replica set
  ```

  ```
  pod[$$first=10,offset=0]{*}.replicaset[$$first=1]{*}
    return first 10 pods which has relicaset and each pod return 1 replicaset
  ```

  ```
  replicaset[@count(pod)<3]{*}.pod{*}
    return replicaset which running pods count less than 3
  ```

## QSL Queries and Their DGraph Equivalents
  ```
  qsl: cluster[@name="preprod-west2.cluster.k8s.local"]{@name}

  dgraph: {
            A as var(func: eq(objtype, cluster)) @filter( eq(name,"preprod-west2.cluster.k8s.local") ) @cascade {
          	  count(uid)
            }
          }
          { objects(func: uid(A),first:1000,offset:0) {
          	  name
          	  uid
            }
          }
  ```
  ```
  qsl: cluster[@name="preprod-west2.cluster.k8s.local"]{*}.namespace[@name="opa"|@name="default"]{*}

  dgraph: {
            A as var(func: eq(objtype, cluster)) @filter( eq(name,"preprod-west2.cluster.k8s.local") ) @cascade {
              count(uid)
              ~cluster @filter(eq(objtype, namespace) and eq(name,"opa") ){
                count(uid)
              }
            }
          }
          { objects(func: uid(A),first:1000,offset:0) {
              creationtime
              k8sobj
              objtype
              name
              resourceid
              resourceversion
              uid
              ~cluster @filter(eq(objtype, namespace) and eq(name,"opa") )(first:1000,offset:0){
                name
                resourceid
                labels
                resourceversion
                creationtime
                k8sobj
                objtype
                uid
              }
            }
          }
  ```

  ```
  qsl: cluster[@name="preprod-west2.cluster.k8s.local"]{@name,@creationtime}.deployment[@name="tiller-deploy"]{@name,@strategy}.ReplicaSet[@name="tiller-deploy-8c8c79584"]{*}

  dgraph: {
            A as var(func: eq(objtype, cluster)) @filter( eq(name,"preprod-west2.cluster.k8s.local") ) @cascade {
           	  count(uid)
           	  ~cluster @filter(eq(objtype, deployment) and eq(name,"tiller-deploy") ){
           	    count(uid)
           	    ~owner @filter(eq(objtype, replicaset) and eq(name,"tiller-deploy-8c8c79584") ){
           	      count(uid)
                }
              }
            }
          }
          { objects(func: uid(A),first:1000,offset:0) {
           	name
           	creationtime
           	uid
           	~cluster @filter(eq(objtype, deployment) and eq(name,"tiller-deploy") )(first:1000,offset:0){
              name
           	  strategy
              uid
           	  ~owner @filter(eq(objtype, replicaset) and eq(name,"tiller-deploy-8c8c79584") )(first:1000,offset:0){
           	    podspec
           		resourceid
           		labels
           		name
           		numreplicas
           		resourceversion
           		creationtime
           		k8sobj
           		objtype
           		uid
              }
            }
           }
          }
  ```
  ```
  qsl:  cluster[@name="preprod-west2.cluster.k8s.local"]{*}.pod[@name="calico-node-gcj7s"]{*}

  dgraph: { A as var(func: eq(objtype, cluster)) @filter( eq(name,"preprod-west2.cluster.k8s.local") ) @cascade {
          	count(uid)
          	~cluster @filter(eq(objtype, pod) and eq(name,"calico-node-gcj7s") ){
          	count(uid)
          }
          }
          }
          { objects(func: uid(A),first:1000,offset:0) {
          	creationtime
          	k8sobj
          	objtype
          	name
          	resourceid
          	resourceversion
          	uid
          	~cluster @filter(eq(objtype, pod) and eq(name,"calico-node-gcj7s") )(first:1000,offset:0){
          		k8sobj
          		ownertype
          		creationtime
          		objtype
          		name
          		containers
          		ip
          		resourceid
          		labels
          		starttime
          		phase
          		resourceversion
          		volumes
          		uid
          }
          }
          }
  ```

  ```
  qsl: namespace[@name="default"]{**}
  dgraph:  { A as var(func: eq(objtype, namespace)) @filter( eq(name,"default") ) @cascade {
           	count(uid)
           }
           }
           { objects(func: uid(A),first:1000,offset:0) {
           	expand(_all_){
           		expand(_all_){
           		}
           	}
           }
           }
  ```

  ```
  qsl: cluster[@objtype="Cluster"$$first=2]{*}.namespace[@name="default"$$first=2,offset=2]{*}
  dgraph: { A as var(func: eq(objtype, cluster)) @filter( eq(objtype,"cluster") ) @cascade {
          	count(uid)
          	~cluster @filter(eq(objtype, namespace) and eq(name,"default") ){
          	count(uid)
          }
          }
          }
          { objects(func: uid(A),first: 2) {
          	creationtime
          	k8sobj
          	objtype
          	name
          	resourceid
          	resourceversion
          	uid
          	~cluster @filter(eq(objtype, namespace) and eq(name,"default") )(first: 2,offset: 2){
          		name
          		resourceid
          		labels
          		resourceversion
          		creationtime
          		k8sobj
          		objtype
          		uid
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
            "resourceid": "preprod-west2.cluster.k8s.local:default",
            "resourceversion": "9",
            "uid": "0x714"
        }
    ]
}
```

```
input:
cluster[@name="preprod-west2.cluster.k8s.local"]{*}.namespace[@name="opa"]{*}

response:
200 OK
{
    "objects": [
        {
            "k8sobj": "K8sObj",
            "name": "preprod-west2.cluster.k8s.local",
            "objtype": "Cluster",
            "resourceid": "preprod-west2.cluster.k8s.local",
            "resourceversion": "0",
            "~cluster": [
                {
                    "k8sobj": "K8sObj",
                    "labels": "null",
                    "name": "opa",
                    "objtype": "Namespace",
                    "resourceid": "preprod-west2.cluster.k8s.local:opa",
                    "resourceversion": "772"
                }
            ]
        }
    ]
}
```

```
input:
cluster[@name="preprod-west2.cluster.k8s.local"]{*}.deployment[@name="dgraph-ratel"]{*}.replicaSet[@name="dgraph-ratel-588856bd5b",@numreplicas>="1"]{*}

response:
200 OK
{
    "objects": [
        {
            "k8sobj": "K8sObj",
            "name": "preprod-west2.cluster.k8s.local",
            "objtype": "Cluster",
            "resourceid": "preprod-west2.cluster.k8s.local",
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
                    "resourceid": "preprod-west2.cluster.k8s.local:dgraph-ratel",
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
                    "resourceid": "preprod-west2.cluster.k8s.local:dgraph-ratel",
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
                            "resourceid": "preprod-west2.cluster.k8s.local:dgraph-ratel-588856bd5b",
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
                    "resourceid": "preprod-west2.cluster.k8s.local:dgraph-ratel",
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
cluster[@name="preprod-west2.cluster.k8s.local"]n{*}.namespace[@name="opa"]{*}

response:
400 Bad Request
Malformed Query: cluster[@name="preprod-west2.cluster.k8s.local"]n{*}.namespace[@name="opa"]{*}
```

#### Error Connecting to Dgraph
```
input:
cluster[@name="preprod-west2.cluster.k8s.local"]{*}.notrealobject[@name="something"]{*}

response:
500 Internal Server Error
Failed to connect to dgraph to get metadata
```

#### Invalid Relation
```
input:
cluster[@name="preprod-west2.cluster.k8s.local"]{*}.notrealobject[@name="something"]{*}

response:
400 Bad Request
no relation found between notrealobject and Cluster
```
