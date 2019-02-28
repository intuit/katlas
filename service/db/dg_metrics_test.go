package db

import (
	"github.com/intuit/katlas/service/metrics"
	"github.com/intuit/katlas/service/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	dc := NewDGClient("127.0.0.1:9080")
	dc.CreateSchema(Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
}

func TestMetricsDgraphNumQueries(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumQueries float64

	query := `
                schema {}
    `

	client := NewDGClient("127.0.0.1:9080")
	defer client.Close()

	//GetQueryResult
	prevCounter = util.ReadCounter(metrics.DgraphNumQueries)
	client.GetQueryResult(query)
	expectedDgraphNumQueries = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumQueries)
	assert.Equal(t, expectedDgraphNumQueries, nextCounter-prevCounter, "DgraphNumQueries is not equal to expected.")

	//GetEntity
	prevCounter = util.ReadCounter(metrics.DgraphNumQueries)
	nid := "0x01"
	_, _ = client.GetEntity(nid)
	expectedDgraphNumQueries = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumQueries)
	assert.Equal(t, expectedDgraphNumQueries, nextCounter-prevCounter, "DgraphNumQueries is not equal to expected.")

	//GetSchemaFromDB
	prevCounter = util.ReadCounter(metrics.DgraphNumQueries)
	_, _ = client.GetSchemaFromDB()
	expectedDgraphNumQueries = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumQueries)
	assert.Equal(t, expectedDgraphNumQueries, nextCounter-prevCounter, "DgraphNumQueries is not equal to expected.")

	//ExecuteDgraphQuery
	prevCounter = util.ReadCounter(metrics.DgraphNumQueries)
	_, _ = client.ExecuteDgraphQuery(query)
	expectedDgraphNumQueries = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumQueries)
	assert.Equal(t, expectedDgraphNumQueries, nextCounter-prevCounter, "DgraphNumQueries is not equal to expected.")
}

func TestMetricsDgraphNumQueriesErr(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumQueriesErr float64
	var client *DGClient

	query := `
                schema {}
    `

	//GetQueryResult
	prevCounter = util.ReadCounter(metrics.DgraphNumQueriesErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	client.GetQueryResult(query)
	expectedDgraphNumQueriesErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumQueriesErr)
	assert.Equal(t, expectedDgraphNumQueriesErr, nextCounter-prevCounter, "DgraphNumQueriesErr is not equal to expected.")

	//GetEntity
	prevCounter = util.ReadCounter(metrics.DgraphNumQueriesErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	var nid = "nid"
	_, _ = client.GetEntity(nid)
	expectedDgraphNumQueriesErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumQueriesErr)
	assert.Equal(t, expectedDgraphNumQueriesErr, nextCounter-prevCounter, "DgraphNumQueriesErr is not equal to expected.")

	//GetSchemaFromDB
	prevCounter = util.ReadCounter(metrics.DgraphNumQueriesErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	_, _ = client.GetSchemaFromDB()
	expectedDgraphNumQueriesErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumQueriesErr)
	assert.Equal(t, expectedDgraphNumQueriesErr, nextCounter-prevCounter, "DgraphNumQueriesErr is not equal to expected.")

	//ExecuteDgraphQuery
	prevCounter = util.ReadCounter(metrics.DgraphNumQueriesErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	_, _ = client.ExecuteDgraphQuery(query)
	expectedDgraphNumQueriesErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumQueriesErr)
	assert.Equal(t, expectedDgraphNumQueriesErr, nextCounter-prevCounter, "DgraphNumQueriesErr is not equal to expected.")
}

func TestMetricsDgraphNumMutations(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumMutations float64

	client := NewDGClient("127.0.0.1:9080")
	defer client.Close()

	//DeleteEntity
	prevCounter = util.ReadCounter(metrics.DgraphNumMutations)
	uid := "0x12345"
	client.DeleteEntity(uid)
	expectedDgraphNumMutations = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumMutations)
	assert.Equal(t, expectedDgraphNumMutations, nextCounter-prevCounter, "DgraphNumMutations is not equal to expected.")

	//CreateEntity
	prevCounter = util.ReadCounter(metrics.DgraphNumMutations)
	node1 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid1, _ := client.CreateEntity("k8snode", node1)
	defer client.DeleteEntity(nid1)
	expectedDgraphNumMutations = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumMutations)
	assert.Equal(t, expectedDgraphNumMutations, nextCounter-prevCounter, "DgraphNumMutations is not equal to expected.")

	//CreateOrDeleteEdge
	prevCounter = util.ReadCounter(metrics.DgraphNumMutations)
	var pid, nodeid string
	client.CreateOrDeleteEdge("k8spod", pid, "k8snode", nodeid, "runsOn", 0)
	expectedDgraphNumMutations = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumMutations)
	assert.Equal(t, expectedDgraphNumMutations, nextCounter-prevCounter, "DgraphNumMutations is not equal to expected.")

	//UpdateEntity
	node2 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid2, _ := client.CreateEntity("k8snode", node2)
	defer client.DeleteEntity(nid2)
	prevCounter = util.ReadCounter(metrics.DgraphNumMutations)
	update := make(map[string]interface{})
	client.UpdateEntity(nid2, update)
	expectedDgraphNumMutations = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumMutations)
	assert.Equal(t, expectedDgraphNumMutations, nextCounter-prevCounter, "DgraphNumMutations is not equal to expected.")
}

func TestMetricsDgraphNumMutationsErr(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumMutationsErr float64
	var client *DGClient

	//DeleteEntity
	prevCounter = util.ReadCounter(metrics.DgraphNumMutationsErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	uid := "0x12345"
	client.DeleteEntity(uid)
	expectedDgraphNumMutationsErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumMutationsErr)
	assert.Equal(t, expectedDgraphNumMutationsErr, nextCounter-prevCounter, "DgraphNumMutationsErr is not equal to expected.")

	//CreateEntity
	prevCounter = util.ReadCounter(metrics.DgraphNumMutationsErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	node1 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid1, _ := client.CreateEntity("k8snode", node1)
	defer client.DeleteEntity(nid1)
	expectedDgraphNumMutationsErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumMutationsErr)
	assert.Equal(t, expectedDgraphNumMutationsErr, nextCounter-prevCounter, "DgraphNumMutationsErr is not equal to expected.")

	//CreateOrDeleteEdge
	prevCounter = util.ReadCounter(metrics.DgraphNumMutationsErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	var pid, nodeid string
	client.CreateOrDeleteEdge("k8spod", pid, "k8snode", nodeid, "runsOn", 0)
	expectedDgraphNumMutationsErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumMutationsErr)
	assert.Equal(t, expectedDgraphNumMutationsErr, nextCounter-prevCounter, "DgraphNumMutationsErr is not equal to expected.")

	//UpdateEntity
	node2 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid2, _ := client.CreateEntity("k8snode", node2)
	defer client.DeleteEntity(nid2)
	prevCounter = util.ReadCounter(metrics.DgraphNumMutationsErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	update := make(map[string]interface{})
	client.UpdateEntity(nid2, update)
	expectedDgraphNumMutationsErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumMutationsErr)
	assert.Equal(t, expectedDgraphNumMutationsErr, nextCounter-prevCounter, "DgraphNumMutationsErr is not equal to expected.")
}

func TestMetricsDgraphNumCreateEntityErr(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumCreateEntityErr float64
	var client *DGClient

	//CreateEntity
	prevCounter = util.ReadCounter(metrics.DgraphNumCreateEntityErr)
	client = NewDGClient("127.0.0.1:9080")
	client.Close()
	node1 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid1, _ := client.CreateEntity("k8snode", node1)
	defer client.DeleteEntity(nid1)

	expectedDgraphNumCreateEntityErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumCreateEntityErr)
	assert.Equal(t, expectedDgraphNumCreateEntityErr, nextCounter-prevCounter, "DgraphNumCreateEntityErr is not equal to expected.")
}

func TestMetricsDgraphNumGetEntityErr(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumGetEntityErr float64
	var client *DGClient

	//CreateEntity
	client = NewDGClient("127.0.0.1:9080")
	node1 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid1, _ := client.CreateEntity("k8snode", node1)
	defer client.DeleteEntity(nid1)

	//Get Entity
	prevCounter = util.ReadCounter(metrics.DgraphNumGetEntityErr)
	var nid string
	client.Close()
	_, _ = client.GetEntity(nid)
	expectedDgraphNumGetEntityErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumGetEntityErr)
	assert.Equal(t, expectedDgraphNumGetEntityErr, nextCounter-prevCounter, "DgraphNumGetEntityErr is not equal to expected.")
}

func TestMetricsDgraphNumUpdateEntityErr(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumUpdateEntityErr float64
	var client *DGClient

	client = NewDGClient("127.0.0.1:9080")
	node1 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid1, _ := client.CreateEntity("k8snode", node1)
	defer client.DeleteEntity(nid1)

	//Update Entity
	prevCounter = util.ReadCounter(metrics.DgraphNumUpdateEntityErr)
	update := make(map[string]interface{})
	client.Close()
	client.UpdateEntity(nid1, update)
	expectedDgraphNumUpdateEntityErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumUpdateEntityErr)
	assert.Equal(t, expectedDgraphNumUpdateEntityErr, nextCounter-prevCounter, "DgraphNumUpdateEntityErr is not equal to expected.")
}

func TestMetricsDgraphNumDeleteEntityErr(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumDeleteEntityErr float64
	var client *DGClient

	client = NewDGClient("127.0.0.1:9080")
	node1 := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid1, _ := client.CreateEntity("k8snode", node1)

	prevCounter = util.ReadCounter(metrics.DgraphNumDeleteEntityErr)
	client.Close()
	client.DeleteEntity(nid1)
	expectedDgraphNumDeleteEntityErr = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumDeleteEntityErr)
	assert.Equal(t, expectedDgraphNumDeleteEntityErr, nextCounter-prevCounter, "DgraphNumDeleteEntityErr is not equal to expected.")
}

func cleanup(nid string) {
	client := NewDGClient("127.0.0.1:9080")
	client.DeleteEntity(nid)
}
