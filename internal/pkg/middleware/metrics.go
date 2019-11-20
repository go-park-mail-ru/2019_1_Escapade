package middleware

import "github.com/prometheus/client_golang/prometheus"

var (
	// Hits hits
	Hits *prometheus.CounterVec
	// Users - registered users
	Users *prometheus.GaugeVec
)

// Init prometheus metrics variables
func Init() {
	var (
		subsystem = "api"
	)
	Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "api",
		Subsystem: subsystem,
	}, []string{"status", "path", "method"})

	prometheus.MustRegister(Hits)
}
