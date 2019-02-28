package apis

import (
	//"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	lru "github.com/hashicorp/golang-lru"
	"github.com/intuit/katlas/service/db"
	"github.com/intuit/katlas/service/metrics"
	"github.com/intuit/katlas/service/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	dc := db.NewDGClient("127.0.0.1:9080")
	dc.CreateSchema(db.Schema{Predicate: "name", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "objtype", Type: "string", Index: true, Tokenizer: []string{"term"}})
	dc.CreateSchema(db.Schema{Predicate: "resourceid", Type: "string", Index: true, Tokenizer: []string{"term"}})
}

func TestMetricsDgraphNumKeywordQueries(t *testing.T) {

	prevCounter := util.ReadCounter(metrics.DgraphNumKeywordQueries)

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
	nextCounter := util.ReadCounter(metrics.DgraphNumKeywordQueries)
	assert.Equal(t, expectedDgraphNumKeywordQueries, nextCounter-prevCounter, "DgraphNumKeywordQueries is not equal to expected.")
}

func TestMetricsDgraphNumKeyValueQueries(t *testing.T) {

	prevCounter := util.ReadCounter(metrics.DgraphNumKeyValueQueries)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewQueryService(dc)
	//Create query map
	m := map[string][]string{
		"name": {"pod01"},
	}
	_, _ = s.GetQueryResult(m)

	expectedDgraphNumKeyValueQueries := 1.0
	nextCounter := util.ReadCounter(metrics.DgraphNumKeyValueQueries)
	assert.Equal(t, expectedDgraphNumKeyValueQueries, nextCounter-prevCounter, "DgraphNumKeyValueQueries is not equal to expected.")
}

func TestMetricsDgraphNumQSL(t *testing.T) {

	prevCounter := util.ReadCounter(metrics.DgraphNumQSL)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	qslSvc := NewQSLService(dc)
	q := `
                cluster{*}
        `
	_, _ = qslSvc.CreateDgraphQuery(q, false)

	expectedDgraphNumQSL := 1.0
	nextCounter := util.ReadCounter(metrics.DgraphNumQSL)
	assert.Equal(t, expectedDgraphNumQSL, nextCounter-prevCounter, "DgraphNumQSL is not equal to expected.")
}

func TestMetricsDgraphCreateEntity(t *testing.T) {

	prevCounter := util.ReadCounter(metrics.DgraphNumCreateEntity)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	node := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid, _ := s.CreateEntity("k8snode", node)
	defer s.DeleteEntity(nid)
	expectedDgraphNumCreateEntity := 1.0
	nextCounter := util.ReadCounter(metrics.DgraphNumCreateEntity)
	assert.Equal(t, expectedDgraphNumCreateEntity, nextCounter-prevCounter, "DgraphNumCreateEntity is not equal to expected.")
}

func TestMetricsDgraphGetEntity(t *testing.T) {

	prevCounter := util.ReadCounter(metrics.DgraphNumGetEntity)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	node := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid, _ := s.CreateEntity("k8snode", node)
	defer s.DeleteEntity(nid)
	_, _ = s.GetEntity(nid)

	expectedDgraphNumGetEntity := 1.0
	nextCounter := util.ReadCounter(metrics.DgraphNumGetEntity)
	assert.Equal(t, expectedDgraphNumGetEntity, nextCounter-prevCounter, "DgraphNumGetEntity is not equal to expected.")
}

func TestMetricsDgraphUpdateEntity(t *testing.T) {

	prevCounter := util.ReadCounter(metrics.DgraphNumUpdateEntity)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	node := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid, _ := s.CreateEntity("k8snode", node)
	defer s.DeleteEntity(nid)
	update := make(map[string]interface{})
	s.UpdateEntity(nid, update)

	expectedDgraphNumUpdateEntity := 1.0
	nextCounter := util.ReadCounter(metrics.DgraphNumUpdateEntity)
	assert.Equal(t, expectedDgraphNumUpdateEntity, nextCounter-prevCounter, "DgraphNumUpdateEntity is not equal to expected.")
}

func TestMetricsDgraphDeleteEntity(t *testing.T) {

	var prevCounter, nextCounter, expectedDgraphNumDeleteEntity float64

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)

	//DeleteEntity
	prevCounter = util.ReadCounter(metrics.DgraphNumDeleteEntity)
	node := map[string]interface{}{
		"objtype": "k8snode",
		"name":    "test-node-metrics",
	}
	nid, _ := s.CreateEntity("k8snode", node)
	s.DeleteEntity(nid)
	expectedDgraphNumDeleteEntity = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumDeleteEntity)
	assert.Equal(t, expectedDgraphNumDeleteEntity, nextCounter-prevCounter, "DgraphNumDeleteEntity is not equal to expected.")

	//DeleteEntityByResourceID
	prevCounter = util.ReadCounter(metrics.DgraphNumDeleteEntity)
	// create node
	nodenew := map[string]interface{}{
		"objtype":    "k8snode",
		"name":       "test-node-metrics",
		"resourceid": "noderid",
	}
	s.CreateEntity("k8snode", nodenew)
	_ = s.DeleteEntityByResourceID("k8snode", "noderid")
	expectedDgraphNumDeleteEntity = 1.0
	nextCounter = util.ReadCounter(metrics.DgraphNumDeleteEntity)
	assert.Equal(t, expectedDgraphNumDeleteEntity, nextCounter-prevCounter, "DgraphNumDeleteEntity is not equal to expected.")
}

func TestMetricsDgraphUpdateEdge(t *testing.T) {

	prevCounter := util.ReadCounter(metrics.DgraphNumUpdateEdge)

	dc := db.NewDGClient("127.0.0.1:9080")
	defer dc.Close()
	s := NewEntityService(dc)
	var pid, nid string
	s.CreateOrDeleteEdge("k8spod", pid, "k8snode", nid, "runsOn", 0)

	expectedDgraphNumUpdateEdge := 1.0
	nextCounter := util.ReadCounter(metrics.DgraphNumUpdateEdge)
	assert.Equal(t, expectedDgraphNumUpdateEdge, nextCounter-prevCounter, "DgraphNumUpdateEdge is not equal to expected.")
}
