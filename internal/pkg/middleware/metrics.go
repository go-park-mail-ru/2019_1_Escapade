package middleware

import "github.com/prometheus/client_golang/prometheus"

var (
	Hits *prometheus.CounterVec
	Users *prometheus.GaugeVec
)

// Init prometheus metrics variables
func Init() {
	var (
		subsystem = "api"
	)
	Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "requests",
		Subsystem: subsystem,
	}, []string{"ip", "status", "path", "method"})

	Users = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "users",
		Subsystem: subsystem,
	}, []string{"ip", "path", "method"})

	prometheus.MustRegister(Hits, Users)
}
