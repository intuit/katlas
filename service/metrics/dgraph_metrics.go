package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//Track Dgraph and Api Service related Metrics

var (
	//DgraphNumQueries ...The total number of Queries processed by Dgraph for Katlas Service
	DgraphNumQueries = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_queries",
		Help: "The total number of Queries processed by Dgraph for Katlas Service",
	})
	//DgraphNumMutations ...The total number of Mutations processed by Dgraph for Katlas Service
	DgraphNumMutations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_mutations",
		Help: "The total number of Mutations processed by Dgraph for Katlas Service",
	})
	//DgraphNumQueriesErr ...The total number of Queries Errored in Dgraph for Katlas Service
	DgraphNumQueriesErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_queries_err",
		Help: "The total number of Queries Errored in Dgraph for Katlas Service",
	})
	//DgraphNumMutationsErr ...The total number of Mutations errored in Dgraph for Katlas Service
	DgraphNumMutationsErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_mutations_err",
		Help: "The total number of Mutations errored in Dgraph for Katlas Service",
	})
	//DgraphNumKeywordQueries ...The total number of Keyword Queries processed by Katlas Query Service
	DgraphNumKeywordQueries = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_num_keyword_queries",
		Help: "The total number of Keyword Queries processed by Katlas Query Service",
	})
	//DgraphNumKeywordQueriesErr ...The total number of Keyword Queries errored in Katlas Query Service
	DgraphNumKeywordQueriesErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_num_keyword_queries_err",
		Help: "The total number of Keyword Queries errored in Katlas Query Service",
	})
	//DgraphNumKeyValueQueries ...The total number of KeyValue Queries processed by Katlas Query Service
	DgraphNumKeyValueQueries = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_num_keyvalue_queries",
		Help: "The total number of KeyValue Queries processed by Katlas Query Service",
	})
	//DgraphNumKeyValueQueriesErr ...The total number of KeyValue Queries errored in Katlas Query Service
	DgraphNumKeyValueQueriesErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_num_keyvalue_queries_err",
		Help: "The total number of KeyValue Queries errored in Katlas Query Service",
	})
	//DgraphNumQSL ...The total number of QSL processed by for Katlas QSL Service
	DgraphNumQSL = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_num_qsl",
		Help: "The total number of QSL processed by Katlas QSL Service",
	})
	//DgraphNumQSLErr ...The total number of QSL errored in Katlas QSL Service
	DgraphNumQSLErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_num_qsl_err",
		Help: "The total number of QSL errored in Katlas QSL Service",
	})
	//DgraphNumCreateEntity ...The total number of Create Entity Requests processed by Katlas Entity Service
	DgraphNumCreateEntity = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_create_entity",
		Help: "The total number of Create Entity Requests processed by Katlas Entity Service",
	})
	//DgraphNumCreateEntityErr ...The total number of Create Entity Requests errored in Katlas Entity Service
	DgraphNumCreateEntityErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_create_entity_err",
		Help: "The total number of Create Entity Requests errored in Katlas Entity Service",
	})
	//DgraphNumUpdateEntity ...The total number of Update Entity Requests processed by Katlas Entity Service
	DgraphNumUpdateEntity = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_update_entity",
		Help: "The total number of Update Entity Requests processed by Katlas Entity Service",
	})
	//DgraphNumUpdateEntityErr ...The total number of Update Entity Requests errored in Katlas Entity Service
	DgraphNumUpdateEntityErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_update_entity_err",
		Help: "The total number of Update Entity Requests errored in Katlas Entity Service",
	})
	//DgraphNumDeleteEntity ...The total number of Delete Entity Requests processed by Katlas Entity Service
	DgraphNumDeleteEntity = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_delete_entity",
		Help: "The total number of Delete Entity Requests processed by Katlas Entity Service",
	})
	//DgraphNumDeleteEntityErr ...The total number of Delete Entity Requests errored in Katlas Entity Service
	DgraphNumDeleteEntityErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_delete_entity_err",
		Help: "The total number of Delete Entity Requests errored in Katlas Entity Service",
	})
	//DgraphNumGetEntity ...The total number of Get Entity Requests processed by Katlas Entity Service
	DgraphNumGetEntity = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_get_entity",
		Help: "The total number of Get Entity Requests processed by Katlas Entity Service",
	})
	//DgraphNumGetEntityErr ...The total number of Get Entity Requests errored in Katlas Entity Service
	DgraphNumGetEntityErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_get_entity_err",
		Help: "The total number of Get Entity Requests errored in Katlas Entity Service",
	})
	//DgraphNumUpdateEdge ...The total number of Create/Delete Edge Requests processed by Katlas Entity Service
	DgraphNumUpdateEdge = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_dgraph_num_update_edge",
		Help: "The total number of Create/Delete Edge Requests processed by Katlas Entity Service",
	})
	//DgraphCreateEntityLatencyHistogram ...latency metric for Create Entity
	DgraphCreateEntityLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "katlas_dgraph_create_entity_latency",
		Help:    "Time take to handle create entity queries by the Katlas Entity Service",
		Buckets: prometheus.ExponentialBuckets(0.0010, 2, 15),
	}, []string{"code"})
	//DgraphUpdateEntityLatencyHistogram ...latency metric for Update Entity
	DgraphUpdateEntityLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "katlas_dgraph_update_entity_latency",
		Help:    "Time take to handle update entity queries by the Katlas Entity Service",
		Buckets: prometheus.ExponentialBuckets(0.0001, 2, 15),
	}, []string{"code"})
	//DgraphDeleteEntityLatencyHistogram ...latency metric for Delete Entity
	DgraphDeleteEntityLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "katlas_dgraph_delete_entity_latency",
		Help:    "Time take to handle delete entity queries by the Katlas Entity Service",
		Buckets: prometheus.ExponentialBuckets(0.0010, 2, 15),
	}, []string{"code"})
	//DgraphGetEntityLatencyHistogram ...latency metric for Get Entity
	DgraphGetEntityLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "katlas_dgraph_get_entity_latency",
		Help:    "Time take to handle get entity queries by the Katlas Entity Service",
		Buckets: prometheus.ExponentialBuckets(0.0010, 2, 15),
	}, []string{"code"})
)
