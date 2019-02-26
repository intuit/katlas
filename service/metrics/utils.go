package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

//RegisterHistogramMetrics ...Register histogram with prometheus
func RegisterHistogramMetrics() {
	prometheus.MustRegister(KatlasQueryLatencyHistogram)
	prometheus.MustRegister(DgraphCreateEntityLatencyHistogram)
	prometheus.MustRegister(DgraphUpdateEntityLatencyHistogram)
	prometheus.MustRegister(DgraphDeleteEntityLatencyHistogram)
	prometheus.MustRegister(DgraphGetEntityLatencyHistogram)
}

//ReadCounter ...Extract float64 Value from the prometheus Counter metric
func ReadCounter(m prometheus.Counter) float64 {
	pb := &dto.Metric{}
	m.Write(pb)
	return pb.GetCounter().GetValue()
}
