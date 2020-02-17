package milogger

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
)

// Logger is implementation of Middleware interface(package infrastructure)
// log every request
type Logger struct {
	server.ServerAddr
	logger infrastructure.Logger
}

// New instance of Logger
func New(logger infrastructure.Logger) *Logger {
	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}
	return &Logger{logger: logger}
}

func (mw *Logger) Func(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			ip := mw.IP(nil)
			mw.logger.Println("listen for you on " + ip)
			next.ServeHTTP(rw, r)
		},
	)
}
