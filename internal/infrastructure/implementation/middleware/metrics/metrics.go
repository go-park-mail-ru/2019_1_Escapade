package mimetrics

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// Metrics is implementation of Middleware interface(package infrastructure)
// record metrics
type Metrics struct {
	server.ServerAddr

	metrics infrastructure.Metrics
	subnet  string
}

// New instance of Metrics
func New(
	metrics infrastructure.Metrics,
	subnet string,
) (*Metrics, error) {
	if metrics == nil {
		return nil, errors.New(ErrNoMetrics)
	}
	return &Metrics{
		metrics: metrics,
		subnet:  subnet,
	}, nil
}

type respWriterStatusCode struct {
	http.ResponseWriter
	status int
}

func (rw *respWriterStatusCode) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (mw *Metrics) Func(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			ip := mw.IP(&mw.subnet)
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
