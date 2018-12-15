---
description: >-
  Describes the core components to help get started with the Installation
  Section.
---

# Core Components

#### Controller

The controller is responsible for discovery of Kubernetes assets in Kubernetes Clusters. The controller must reside in the Kubernetes cluster for which it will collect asset data. For details on Controller design, please refer [Design Concepts](design-concepts.md)

#### Rest Service

The Rest Service exposes APIs that can be used to get details about Kubernetes entities, run queries to help diagnose issues in kubernetes clusters. The rest service must reside in the Cutlass Kubernetes Cluster. For details on Rest Service Calls, please refer [Rest APIs](rest-apis.md) 

#### Web Application

The Web Application exposes UI search capability to search clusters based on several criteria. The Web application must be hosted on the Cutlass Kubernetes Cluster. For details on usage, please refer [Cutlass UI]()

#### Database

Dgraph is used as the graph database. To know more about our motivation to choose Dgraph, please refer [Design Concepts](design-concepts.md)



