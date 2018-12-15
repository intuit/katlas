# Motivation

In a Kubernetes ecosystem, with multiple clusters deployed, it is not only time consuming, but tedious for operators/developers to debug issues. Application deployment data are stored in various dispersed systems and difficult to navigate, audit and diagnose. K-Atlas solves this by providing a holistic view and query interface.

Additionally, there is no graphical view of the topology, making it hard to visualize the system. K-Atlas solves this by providing a near real-time graphical view of the system.

## **User Cases**

#### **View and Report on Cluster Consistency**

* Evaluate the configuration of multiple K8s clusters \(potentially in multiple regions\)
* Compare against baseline and each other, and provide a diff

#### **Application Centric View**

* Present a global view of an application - across clusters and regions
* What versions of the application are running? Where?
* What are the key objects powering the app? What is their current status? How are they changing?
* Monitor application fleet on a single pane of glass - how many resources are my application consuming? In what regions?
* Diagnose and narrow down application problems
  * How many pods are backing a dns domain x.y.abc.com?
  * What is the application deployment that has pod IP a.b.c.d?

#### **Policy Compliance and Enforcement**

* Set policies and query for compliance
* Send out alerts when policies are violated and aid in enforcement
* Examples of policies include:
  * Does application deployment meet minimum DR requirement?
  * Run only approved protocols on approved ports \(HTTPS on 443\)

#### **Change Management**

* Playback changes that happened on the fleet
* Provide deep visibility into any change to the fleet









