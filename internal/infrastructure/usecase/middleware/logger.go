package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
)

type Logger struct {
	logger infrastructure.LoggerI
}

func NewLogger(logger infrastructure.LoggerI) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (mw *Logger) Func(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			ip := server.GetIP(nil)
			mw.logger.Println("listen for you on " + ip)
			next.ServeHTTP(rw, r)
		},
	)
}
