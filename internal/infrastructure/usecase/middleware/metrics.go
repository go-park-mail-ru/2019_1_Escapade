package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

type Metrics struct {
	metrics infrastructure.MetricsI
	subnet  string
}

func NewMetrics(
	metrics infrastructure.MetricsI,
	subnet string,
) *Metrics {
	return &Metrics{
		metrics: metrics,
		subnet:  subnet,
	}
}

type respWriterStatusCode struct {
	http.ResponseWriter
	status int
}

func (rw *respWriterStatusCode) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Metrics record metrics
func (mw *Metrics) Func(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			ip := server.GetIP(&mw.subnet)
			goodRW := &respWriterStatusCode{rw, 200}
			mw.metrics.UsersInc(ip, r.URL.Path, r.Method)
			next.ServeHTTP(goodRW, r)
			mw.metrics.HitsInc(
				ip,
				utils.String(goodRW.status),
				r.URL.Path,
				r.Method,
			)
			mw.metrics.UsersDec(ip, r.URL.Path, r.Method)
		},
	)
}
