package prometheus

import "github.com/prometheus/client_golang/prometheus"

type Prometheus struct {
	Hits  *prometheus.CounterVec
	Users *prometheus.GaugeVec
}

// Init prometheus metrics variables
func New() *Prometheus {
	var (
		pr        Prometheus
		subsystem = "api"
	)
	pr.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "requests",
		Subsystem: subsystem,
	}, []string{"ip", "status", "path", "method"})

	pr.Users = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "users",
		Subsystem: subsystem,
	}, []string{"ip", "path", "method"})

	prometheus.MustRegister(pr.Hits, pr.Users)
	return &pr
}

func (pr *Prometheus) HitsInc(ip, status, path, method string) {
	pr.Hits.WithLabelValues(ip, status, path, method).Inc()
}

func (pr *Prometheus) UsersInc(ip, path, method string) {
	pr.Users.WithLabelValues(ip, path, method).Inc()
}

func (pr *Prometheus) UsersDec(ip, path, method string) {
	pr.Users.WithLabelValues(ip, path, method).Dec()
}
