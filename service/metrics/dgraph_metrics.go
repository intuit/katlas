package metrics

import (
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
        //DgraphNumQueries ...The total number of Queries processed by Dgraph for Katlas Service
        DgraphNumQueries = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_queries",
                Help: "The total number of Queries processed by Dgraph for Katlas Service",
        })
        //DgraphNumMutations ...The total number of Mutations processed by Dgraph for Katlas Service
        DgraphNumMutations = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_mutations",
                Help: "The total number of Mutations processed by Dgraph for Katlas Service",
        })
        //DgraphNumQueriesErr ...The total number of Queries Errored in Dgraph for Katlas Service
        DgraphNumQueriesErr = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_queries_err",
                Help: "The total number of Queries Errored in Dgraph for Katlas Service",
        })
        //DgraphNumMutationsErr ...The total number of Mutations errored in Dgraph for Katlas Service
        DgraphNumMutationsErr = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_mutations_err",
                Help: "The total number of Mutations errored in Dgraph for Katlas Service",
        })
        //DgraphNumKeywordQueries ...The total number of Keyword Queries processed by Katlas Query Service
        DgraphNumKeywordQueries = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_keyword_queries",
                Help: "The total number of Keyword Queries processed by Katlas Query Service",
        })
        //DgraphNumKeyValueQueries ...The total number of KeyValue Queries processed by Katlas Query Service
        DgraphNumKeyValueQueries = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_keyvalue_queries",
                Help: "The total number of KeyValue Queries processed by Katlas Query Service",
        })
        //DgraphNumQSL ...The total number of QSL processed by for Katlas QSL Service
        DgraphNumQSL = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_qsl",
                Help: "The total number of QSL processed by for Katlas QSL Service",
        })
        //DgraphNumCreateEntity ...The total number of Create Entity Requests processed by Katlas Entity Service
        DgraphNumCreateEntity = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_create_entity",
                Help: "The total number of Create Entity Requests processed by Katlas Entity Service",
        })
        //DgraphNumUpdateEntity ...The total number of Update Entity Requests processed by Katlas Entity Service
        DgraphNumUpdateEntity = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_update_entity",
                Help: "The total number of Update Entity Requests processed by Katlas Entity Service",
        })
        //DgraphNumDeleteEntity ...The total number of Delete Entity Requests processed by Katlas Entity Service
        DgraphNumDeleteEntity = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_delete_entity",
                Help: "The total number of Delete Entity Requests processed by Katlas Entity Service",
        })
        //DgraphNumGetEntity ...The total number of Get Entity Requests processed by Katlas Entity Service
        DgraphNumGetEntity = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_get_entity",
                Help: "The total number of Get Entity Requests processed by Katlas Entity Service",
        })
        //DgraphNumUpdateEdge ...The total number of Create/Delete Edge Requests processed by Katlas Entity Service
        DgraphNumUpdateEdge = promauto.NewCounter(prometheus.CounterOpts{
                Name: "katlas_service_dgraph_num_update_edge",
                Help: "The total number of Create/Delete Edge Requests processed by Katlas Entity Service",
        })
)
