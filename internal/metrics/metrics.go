package metrics

import "github.com/prometheus/client_golang/prometheus"

var Hits *prometheus.CounterVec

func InitHitsMetric(service string) {
	Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "Hits",
		Help: service,
	}, []string{"statsus", "path", "method"})
}
