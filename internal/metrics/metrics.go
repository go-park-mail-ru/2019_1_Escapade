package metrics

import "github.com/prometheus/client_golang/prometheus"

var Hits *prometheus.CounterVec


func GetHitsMetric(service string){
	Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "Hits"
	}, []strings{"statsus", "path", "method"})
}
