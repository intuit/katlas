# K-Atlas

[![Build Status](https://api.travis-ci.org/intuit/katlas.svg?branch=master)](https://www.travis-ci.org/intuit/katlas)
[![codecov](https://codecov.io/gh/intuit/katlas/branch/master/graph/badge.svg)](https://codecov.io/gh/intuit/katlas)
[![Go Report Card](https://goreportcard.com/badge/github.com/intuit/katlas)](https://goreportcard.com/report/github.com/intuit/katlas)
[![Slack Chat](https://img.shields.io/badge/slack-live-orange.svg)](https://katlasio.slack.com/)


![](website/assets/images/katlas-logo-blue-300px.png)

## What It Does?

**K-Atlas** \(_pronounced **Cutlass**_\), is a distributed graph based platform to automatically collect, discover, explore and relate multi-cluster Kubernetes resources and metadata. K-Atlas's rich query language allows for simple and efficient exploration and extensibility.

It addresses following problems in a large scale enterprise environment of Kubernetes.

* **Discoverability**
  * Find K8s objects across multiple distributed K8s clusters
  * Real-time view of discovered objects
  * Streaming APIs and UI for both programmatic and human interactions
* **Advanced Exploration**
  * Identify similarities and differences between objects from pods to clusters
  * Correlate different objects by performing advanced join operations
* **Federated Application View**
  * Applications take center stage. K-Atlas provides a unique, application-centric view, with metadata from multiple clusters
  * Single pane of glass view of the entire application - from edge to database, across all clusters, regions etc.
* **Reporting**
  * Provide advanced reporting on compliance, security and other organizational policies
* **Policy Enforcement**
  * Allow for organizational policies to be enforced across the fleet in a consistent manner

Check out more details on [Motivation and Use Cases](docs/motivation.md) that K-Atlas is addressing.

It provides a Web Viewer that can be used to search the Kubernetes cluster data and view graphical results in real time. Click [here to see a demo](https://www.useloom.com/share/eb97aa1054004be197e3ed732223e689).

## Core Components

![](docs/diagram/K-Atlas.png)

#### Collector

The collector is responsible for discovery of Kubernetes assets in Kubernetes Clusters. For details on the  Collector design, please refer [Design Concepts](docs/design-concepts.md)

#### K-Atlas Service

The K-Atlas Service exposes APIs that can be used to get details about Kubernetes entities and run queries to help diagnose issues in Kubernetes clusters. For details , please refer [K-Atlas APIs](docs/rest-apis.md) 

#### K-Atlas Browser

The Web Application exposes UI search capability to search clusters based on several criteria and provide a real time graphical view of entities. For details on usage, please Click [here to see a demo](https://www.useloom.com/share/eb97aa1054004be197e3ed732223e689)

#### Graph Database

Dgraph is used as the graph database. To know more about our motivation to choose Dgraph, please refer [Design Concepts](docs/design-concepts.md)

## Deploying to a Cluster

### Technical Requirements

Make sure you have the following prerequisites:

* A local Go 1.7+ development environment.
* Access to a Kubernetes cluster.

### Setup Steps

How to [Set Up](docs/installation.md).

## Releases

#### Latest version (v0.6)

* [x] QSL query support
* [x] Dynamic search result layout based on QSL
* [x] Graph view based on QSL query required objects & relationships
* [x] Pagination support for both API and UI
* [x] Custom metadata definition via new API

More details about specific K-Atlas features are at [Release Notes](release.md).

## Contributing

We encourage you to get involved with K-Atlas, as users or contributors and help with code reviews.

Read the [contributing guidelines](docs/contributing.md) to learn about building the project, the project structure, and the purpose of each package. 

