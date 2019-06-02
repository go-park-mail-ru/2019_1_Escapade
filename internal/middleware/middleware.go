package middleware

import (
	"fmt"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/cookie"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/cors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"net/http"
)

type respWriterStatusCode struct {
	http.ResponseWriter
	status int
}

func (rw *respWriterStatusCode) WriteHeader(status int) {
	rw.status = status
	fmt.Println("status", status)
	rw.ResponseWriter.WriteHeader(status)
}

// HandleDecorator middleware
type HandleDecorator func(http.HandlerFunc) http.HandlerFunc

// CORS Access-Control-Allow-Origin
func CORS(cc config.CORSConfig, preCORS bool) HandleDecorator {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			origin := cors.GetOrigin(r)
			if !cors.IsAllowed(origin, cc.Origins) {
				place := "middleware/CORS"
				utils.PrintResult(re.ErrorCORS(origin), http.StatusForbidden, place)
				rw.WriteHeader(http.StatusForbidden)
				return
			}
			//
			cors.SetCORS(rw, cc, origin)

			if preCORS {
				rw.WriteHeader(http.StatusOK)
			} else {
				next(rw, r)
			}
			return
		}
	}
}

// Auth Check cookie exists
func Auth(cc config.SessionConfig) HandleDecorator {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			if _, err := cookie.GetSessionCookie(r, cc); err != nil {
				const place = "middleware/Auth"
				utils.PrintResult(err, http.StatusUnauthorized, place)
				rw.WriteHeader(http.StatusUnauthorized)
				utils.SendErrorJSON(rw, re.ErrorNoCookie(), place)
				return
			}

			next(rw, r)
		}
	}
}

//Recover catch panic
func Recover(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		defer utils.CatchPanic("middleware.go Recover()")

		next(rw, r)
	}
}

func Metrics(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		goodRW := &respWriterStatusCode{rw, 200}
		next(goodRW, r)
		fmt.Println(goodRW.status)
		metrics.Hits.WithLabelValues(strconv.Itoa(goodRW.status), r.URL.Path, r.Method).Inc()
	}
}

// ApplyMiddleware apply middleware
func ApplyMiddleware(handler http.HandlerFunc,
	decorators ...HandleDecorator) http.HandlerFunc {
	handler = Recover(handler)
	fmt.Println("caaaatchssss")
	for _, m := range decorators {
		handler = m(handler)
	}
	handler = Metrics(handler)
	return handler
}
