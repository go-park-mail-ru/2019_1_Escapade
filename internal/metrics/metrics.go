package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	Hits      *prometheus.CounterVec
	Rooms     prometheus.Gauge
	FreeRooms prometheus.Gauge

	WaitingPlayers prometheus.Gauge
	Players        *prometheus.GaugeVec
)

func InitHitsMetric(service string) {
	Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "Hits",
		Namespace: service,
	}, []string{"statsus", "path", "method"})
}

func InitRoomMetric(service string) {
	Rooms = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "Rooms",
		Namespace: service,
		Help:      "Number of active Rooms",
	})
	FreeRooms = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "FreeRooms",
		Namespace: service,
		Help:      "Number of free Rooms",
	})
}

func InitPlayersMetric(service string) {
	Players = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "Players",
		Namespace: service,
		Help:      "Active Players by Rooms",
	}, []string{"room", "player"})

	WaitingPlayers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "WaitingPlayers",
		Namespace: service,
		Help:      "Number of waiting users",
	})
}
