package middleware

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/cookie"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/cors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

type respWriterStatusCode struct {
	http.ResponseWriter
	status int
}

func (rw *respWriterStatusCode) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// CORS Access-Control-Allow-Origin
func CORS(cc config.CORS) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			origin := cors.GetOrigin(r)
			if !cors.IsAllowed(origin, cc.Origins) {
				api.SendResult(rw, api.NewResult(http.StatusForbidden, "cors", nil, re.CORS(origin)))
				utils.Debug(false, "cors no!!!!!!!!!!")
				return
			}

			//utils.Debug(false, "cors allow!")
			cors.SetCORS(rw, cc, origin)
			next.ServeHTTP(rw, r)
			return
		})
	}
}

// Auth Check cookie exists
func Auth(cc config.Cookie, ca config.Auth, client config.AuthClient) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			utils.Debug(false, "auth start")
			var (
				maxLimit = 1
				userID   string
				err      error
			)
			for i := 0; i < maxLimit; i++ {
				utils.Debug(false, "auth check start")
				userID, err = auth.Check(rw, r, cc, ca, client)
				utils.Debug(false, "auth check end")
				if err == nil {
					break
				}
			}
			if err != nil {
				const place = "middleware/Auth"
				api.SendResult(rw, api.NewResult(http.StatusUnauthorized, place, nil, err))
				return
			}
			ctx := context.WithValue(r.Context(), api.ContextUserKey, userID)

			utils.Debug(false, "auth end", userID, api.ContextUserKey)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

//Recover catch panic
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		defer utils.CatchPanic("middleware.go Recover()")

		next.ServeHTTP(rw, r)
	})
}

// Metrics record metrics
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		goodRW := &respWriterStatusCode{rw, 200}
		next.ServeHTTP(goodRW, r)
		metrics.Hits.WithLabelValues(strconv.Itoa(goodRW.status), r.URL.Path, r.Method).Inc()
	})
}
