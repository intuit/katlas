package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//Track Service level Metrics aggregated across all Request Types.

var (
	//KatlasNumReqErr4xx ...The total number of 4xx Requests processed by Katlas Service
	KatlasNumReqErr4xx = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_service_requests_error_4xx",
		Help: "The total number of 4xx requests processed by Katlas Service",
	})

	//KatlasNumReqErr5xx ...The total number of 5xx Requests processed by Katlas Service
	KatlasNumReqErr5xx = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_service_requests_error_5xx",
		Help: "The total number of 5xx requests processed by Katlas Service",
	})

	//KatlasNumReqErr ...The total number of Error Requests processed by Katlas Service
	KatlasNumReqErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_service_requests_errors",
		Help: "The total number of error requests processed by Katlas Service",
	})

	//KatlasNumReq2xx ...The total number of 2xx Requests processed by Katlas Service
	KatlasNumReq2xx = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_service_requests_2xx",
		Help: "The total number of 2xx requests processed by Katlas Service",
	})

	//KatlasNumReqCount ...The total number of Requests processed by Katlas Service
	KatlasNumReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "katlas_service_requests_count",
		Help: "The total number of requests processed by Katlas Service",
	})

	//KatlasQueryLatencyHistogram ...latency metric for external Query Requests(keyword, QSL queries)
	KatlasQueryLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_server_requests_seconds",
		Help:    "Time take to handle external queries by the Katlas Service",
		Buckets: prometheus.ExponentialBuckets(0.0010, 2, 15),
	}, []string{"app", "caller", "exception", "l1", "l2", "method", "status", "uri"})
)
