package util

import (
	"github.com/cenkalti/backoff"
	"github.com/intuit/katlas/service/metrics"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"time"
)

// OptionContext to define options when CRUD entities
type OptionContext struct {
	// is replace field when update
	ReplaceListOrEdge bool
}

// NewBackOff creates an instance of ExponentialBackOff using default values.
func NewBackOff() *backoff.ExponentialBackOff {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     100 * time.Microsecond,
		RandomizationFactor: 0.5,
		Multiplier:          1.5,
		MaxInterval:         500 * time.Microsecond,
		MaxElapsedTime:      5 * time.Minute,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	return b
}

//RegisterHistogramMetrics ...Register histogram with prometheus
func RegisterHistogramMetrics() {
	prometheus.MustRegister(metrics.KatlasQueryLatencyHistogram)
	prometheus.MustRegister(metrics.DgraphCreateEntityLatencyHistogram)
	prometheus.MustRegister(metrics.DgraphUpdateEntityLatencyHistogram)
	prometheus.MustRegister(metrics.DgraphDeleteEntityLatencyHistogram)
	prometheus.MustRegister(metrics.DgraphGetEntityLatencyHistogram)
}

//ReadCounter ...Extract float64 Value from the prometheus Counter metric
func ReadCounter(m prometheus.Counter) float64 {
	pb := &dto.Metric{}
	m.Write(pb)
	return pb.GetCounter().GetValue()
}
