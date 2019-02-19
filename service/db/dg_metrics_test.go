package db

import (
	"github.com/intuit/katlas/service/metrics"
	"github.com/intuit/katlas/service/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricsDgraphNumQueries(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumQueries float64

	query := `
                schema {}
    `

	client := NewDGClient("127.0.0.1:9080")
	defer client.Close()

	//GetQueryResult
	prevCounter = metrics.ReadCounter(metrics.DgraphNumQueries)
	client.GetQueryResult(query)
	expectedDgraphNumQueries = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumQueries)
	assert.Equal(t, expectedDgraphNumQueries, nextCounter-prevCounter, "DgraphNumQueries is not equal to expected.")

	//GetEntity
	prevCounter = metrics.ReadCounter(metrics.DgraphNumQueries)
	nid := "0x01"
	_, _ = client.GetEntity("k8snode", nid)
	expectedDgraphNumQueries = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumQueries)
	assert.Equal(t, expectedDgraphNumQueries, nextCounter-prevCounter, "DgraphNumQueries is not equal to expected.")

	//GetSchemaFromDB
	prevCounter = metrics.ReadCounter(metrics.DgraphNumQueries)
	_, _ = client.GetSchemaFromDB()
	expectedDgraphNumQueries = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumQueries)
	assert.Equal(t, expectedDgraphNumQueries, nextCounter-prevCounter, "DgraphNumQueries is not equal to expected.")

	//ExecuteDgraphQuery
	prevCounter = metrics.ReadCounter(metrics.DgraphNumQueries)
	_, _ = client.ExecuteDgraphQuery(query)
	expectedDgraphNumQueries = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumQueries)
	assert.Equal(t, expectedDgraphNumQueries, nextCounter-prevCounter, "DgraphNumQueries is not equal to expected.")
}

func TestMetricsDgraphNumQueriesErr(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumQueriesErr float64
	var client *DGClient

	query := `
                schema {}
    `

	//GetQueryResult
	prevCounter = metrics.ReadCounter(metrics.DgraphNumQueriesErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	client.GetQueryResult(query)
	expectedDgraphNumQueriesErr = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumQueriesErr)
	assert.Equal(t, expectedDgraphNumQueriesErr, nextCounter-prevCounter, "DgraphNumQueriesErr is not equal to expected.")

	//GetEntity
	prevCounter = metrics.ReadCounter(metrics.DgraphNumQueriesErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	var nid = "nid"
	_, _ = client.GetEntity("k8snode", nid)
	expectedDgraphNumQueriesErr = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumQueriesErr)
	assert.Equal(t, expectedDgraphNumQueriesErr, nextCounter-prevCounter, "DgraphNumQueriesErr is not equal to expected.")

	//GetSchemaFromDB
	prevCounter = metrics.ReadCounter(metrics.DgraphNumQueriesErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	_, _ = client.GetSchemaFromDB()
	expectedDgraphNumQueriesErr = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumQueriesErr)
	assert.Equal(t, expectedDgraphNumQueriesErr, nextCounter-prevCounter, "DgraphNumQueriesErr is not equal to expected.")

	//ExecuteDgraphQuery
	prevCounter = metrics.ReadCounter(metrics.DgraphNumQueriesErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	_, _ = client.ExecuteDgraphQuery(query)
	expectedDgraphNumQueriesErr = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumQueriesErr)
	assert.Equal(t, expectedDgraphNumQueriesErr, nextCounter-prevCounter, "DgraphNumQueriesErr is not equal to expected.")
}

func TestMetricsDgraphNumMutations(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumMutations float64

	client := NewDGClient("127.0.0.1:9080")
	defer client.Close()

	//DeleteEntity
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	uid := "0x12345"
	client.DeleteEntity(uid)
	expectedDgraphNumMutations = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	assert.Equal(t, expectedDgraphNumMutations, nextCounter-prevCounter, "DgraphNumMutations is not equal to expected.")

	//CreateEntity
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	node1 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid1, _ := client.CreateEntity("k8snode", node1)
	defer client.DeleteEntity(nid1)
	expectedDgraphNumMutations = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	assert.Equal(t, expectedDgraphNumMutations, nextCounter-prevCounter, "DgraphNumMutations is not equal to expected.")

	//SetFieldToNull
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	delMap := make(map[string]interface{})
	delMap[util.UID] = "0x12345"
	_ = client.SetFieldToNull(delMap)
	expectedDgraphNumMutations = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	assert.Equal(t, expectedDgraphNumMutations, nextCounter-prevCounter, "DgraphNumMutations is not equal to expected.")

	//CreateOrDeleteEdge
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	var pid, nodeid string
	client.CreateOrDeleteEdge("k8spod", pid, "k8snode", nodeid, "runsOn", 0)
	expectedDgraphNumMutations = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	assert.Equal(t, expectedDgraphNumMutations, nextCounter-prevCounter, "DgraphNumMutations is not equal to expected.")

	//UpdateEntity
	node2 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid2, _ := client.CreateEntity("k8snode", node2)
	defer client.DeleteEntity(nid2)
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	update := make(map[string]interface{})
	client.UpdateEntity("k8snode", nid2, update)
	expectedDgraphNumMutations = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutations)
	assert.Equal(t, expectedDgraphNumMutations, nextCounter-prevCounter, "DgraphNumMutations is not equal to expected.")
}

func TestMetricsDgraphNumMutationsErr(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumMutationsErr float64
	var client *DGClient

	//DeleteEntity
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	uid := "0x12345"
	client.DeleteEntity(uid)
	expectedDgraphNumMutationsErr = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	assert.Equal(t, expectedDgraphNumMutationsErr, nextCounter-prevCounter, "DgraphNumMutationsErr is not equal to expected.")

	//CreateEntity
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	node1 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid1, _ := client.CreateEntity("k8snode", node1)
	defer client.DeleteEntity(nid1)
	expectedDgraphNumMutationsErr = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	assert.Equal(t, expectedDgraphNumMutationsErr, nextCounter-prevCounter, "DgraphNumMutationsErr is not equal to expected.")

	//SetFieldToNull
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	delMap := make(map[string]interface{})
	delMap[util.UID] = "0x12345"
	_ = client.SetFieldToNull(delMap)
	expectedDgraphNumMutationsErr = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	assert.Equal(t, expectedDgraphNumMutationsErr, nextCounter-prevCounter, "DgraphNumMutationsErr is not equal to expected.")

	//CreateOrDeleteEdge
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	var pid, nodeid string
	client.CreateOrDeleteEdge("k8spod", pid, "k8snode", nodeid, "runsOn", 0)
	expectedDgraphNumMutationsErr = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	assert.Equal(t, expectedDgraphNumMutationsErr, nextCounter-prevCounter, "DgraphNumMutationsErr is not equal to expected.")

	//UpdateEntity
	node2 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid2, _ := client.CreateEntity("k8snode", node2)
	defer client.DeleteEntity(nid2)
	prevCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	update := make(map[string]interface{})
	client.UpdateEntity("k8snode", nid2, update)
	expectedDgraphNumMutationsErr = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumMutationsErr)
	assert.Equal(t, expectedDgraphNumMutationsErr, nextCounter-prevCounter, "DgraphNumMutationsErr is not equal to expected.")
}
