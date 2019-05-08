package metrics

import "github.com/prometheus/client_golang/prometheus"

var Hits *prometheus.CounterVec

var Rooms prometheus.Counter

var Players *prometheus.CounterVec

func InitHitsMetric(service string) {
	Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "Hits",
		Namespace: service,
	}, []string{"statsus", "path", "method"})
}

func InitRoomMetric(service string) {
	Rooms = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "Rooms",
		Namespace: service,
		Help:      "Number of active Rooms",
	})
}

func InitPlayersMetric(service string) {
	Players = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "Players",
		Namespace: service,
		Help:      "Active Players by Rooms",
	}, []string{"room", "playes"})
}
