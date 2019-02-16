package apis

import (
	//"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	lru "github.com/hashicorp/golang-lru"
	"github.com/intuit/katlas/service/db"
	"github.com/intuit/katlas/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestMetricsDgraphNumKeywordQueries(t *testing.T) {

	prevCounter := metrics.ReadCounter(metrics.DgraphNumKeywordQueries)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	var err error
	db.LruCache, err = lru.New(5)
	if err != nil {
		log.Errorf("err: %v", err)
	}
	s := NewQueryService(dc)
	//Create query map
	m := map[string][]string{
		"keyword": {"pod"},
	}
	_, _ = s.GetQueryResult(m)

	expectedDgraphNumKeywordQueries := 1.0
	nextCounter := metrics.ReadCounter(metrics.DgraphNumKeywordQueries)
	assert.Equal(t, expectedDgraphNumKeywordQueries, nextCounter-prevCounter, "DgraphNumKeywordQueries is not equal to expected.")
}

func TestMetricsDgraphNumKeyValueQueries(t *testing.T) {

	prevCounter := metrics.ReadCounter(metrics.DgraphNumKeyValueQueries)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewQueryService(dc)
	//Create query map
	m := map[string][]string{
		"name": {"pod01"},
	}
	_, _ = s.GetQueryResult(m)

	expectedDgraphNumKeyValueQueries := 1.0
	nextCounter := metrics.ReadCounter(metrics.DgraphNumKeyValueQueries)
	assert.Equal(t, expectedDgraphNumKeyValueQueries, nextCounter-prevCounter, "DgraphNumKeyValueQueries is not equal to expected.")
}

func TestMetricsDgraphNumQSL(t *testing.T) {

	prevCounter := metrics.ReadCounter(metrics.DgraphNumQSL)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	qslSvc := NewQSLService(dc)
	q := `
                cluster{*}
        `
	_, _ = qslSvc.CreateDgraphQuery(q, false)

	expectedDgraphNumQSL := 1.0
	nextCounter := metrics.ReadCounter(metrics.DgraphNumQSL)
	assert.Equal(t, expectedDgraphNumQSL, nextCounter-prevCounter, "DgraphNumQSL is not equal to expected.")
}

func TestMetricsDgraphCreateEntity(t *testing.T) {

	prevCounter := metrics.ReadCounter(metrics.DgraphNumCreateEntity)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	node := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nids, _ := s.CreateEntity("k8snode", node)
	var nid string
	for _, v := range nids {
		nid = v
		break
	}
	defer s.DeleteEntity(nid)
	expectedDgraphNumCreateEntity := 1.0
	nextCounter := metrics.ReadCounter(metrics.DgraphNumCreateEntity)
	assert.Equal(t, expectedDgraphNumCreateEntity, nextCounter-prevCounter, "DgraphNumCreateEntity is not equal to expected.")
}

func TestMetricsDgraphGetEntity(t *testing.T) {

	prevCounter := metrics.ReadCounter(metrics.DgraphNumGetEntity)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	node := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nids, _ := s.CreateEntity("k8snode", node)
	var nid string
	for _, v := range nids {
		nid = v
		break
	}
	defer s.DeleteEntity(nid)
	_, _ = s.GetEntity("k8snode", nid)

	expectedDgraphNumGetEntity := 1.0
	nextCounter := metrics.ReadCounter(metrics.DgraphNumGetEntity)
	assert.Equal(t, expectedDgraphNumGetEntity, nextCounter-prevCounter, "DgraphNumGetEntity is not equal to expected.")
}

func TestMetricsDgraphUpdateEntity(t *testing.T) {

	prevCounter := metrics.ReadCounter(metrics.DgraphNumUpdateEntity)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	node := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nids, _ := s.CreateEntity("k8snode", node)
	var nid string
	for _, v := range nids {
		nid = v
		break
	}
	defer s.DeleteEntity(nid)
	update := make(map[string]interface{})
	s.UpdateEntity("k8snode", nid, update)

	expectedDgraphNumUpdateEntity := 1.0
	nextCounter := metrics.ReadCounter(metrics.DgraphNumUpdateEntity)
	assert.Equal(t, expectedDgraphNumUpdateEntity, nextCounter-prevCounter, "DgraphNumUpdateEntity is not equal to expected.")
}

func TestMetricsDgraphDeleteEntity(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumDeleteEntity float64

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)

	//DeleteEntity
	prevCounter = metrics.ReadCounter(metrics.DgraphNumDeleteEntity)
	node := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nids, _ := s.CreateEntity("k8snode", node)
	var nid string
	for _, v := range nids {
		nid = v
		break
	}
	s.DeleteEntity(nid)
	expectedDgraphNumDeleteEntity = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumDeleteEntity)
	assert.Equal(t, expectedDgraphNumDeleteEntity, nextCounter-prevCounter, "DgraphNumDeleteEntity is not equal to expected.")

	//DeleteEntityByResourceID
	prevCounter = metrics.ReadCounter(metrics.DgraphNumDeleteEntity)
	// create node
	nodenew := map[string]interface{}{
		"objtype":    "k8snode",
		"name":       "test-node-metrics",
		"resourceid": "noderid",
	}
	s.CreateEntity("k8snode", nodenew)
	_ = s.DeleteEntityByResourceID("k8snode", "noderid")
	expectedDgraphNumDeleteEntity = 1.0
	nextCounter = metrics.ReadCounter(metrics.DgraphNumDeleteEntity)
	assert.Equal(t, expectedDgraphNumDeleteEntity, nextCounter-prevCounter, "DgraphNumDeleteEntity is not equal to expected.")
}

func TestMetricsDgraphUpdateEdge(t *testing.T) {

	prevCounter := metrics.ReadCounter(metrics.DgraphNumUpdateEdge)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	var pid, nid string
	s.CreateOrDeleteEdge("k8spod", pid, "k8snode", nid, "runsOn", 0)

	expectedDgraphNumUpdateEdge := 1.0
	nextCounter := metrics.ReadCounter(metrics.DgraphNumUpdateEdge)
	assert.Equal(t, expectedDgraphNumUpdateEdge, nextCounter-prevCounter, "DgraphNumUpdateEdge is not equal to expected.")
}
