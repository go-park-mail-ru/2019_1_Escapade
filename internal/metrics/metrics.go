package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	Hits      *prometheus.CounterVec
	Rooms     prometheus.Counter
	FreeRooms prometheus.Counter

	WaitingPlayers prometheus.Counter
	Players        *prometheus.CounterVec
)

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
	FreeRooms = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "FreeRooms",
		Namespace: service,
		Help:      "Number of free Rooms",
	})
}

func InitPlayersMetric(service string) {
	Players = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "Players",
		Namespace: service,
		Help:      "Active Players by Rooms",
	}, []string{"room", "player"})

	WaitingPlayers = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "WaitingPlayers",
		Namespace: service,
		Help:      "Number of waiting users",
	})
}
