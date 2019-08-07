package middleware

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/cookie"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/cors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"strconv"

	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type ContextKey string

const ContextUserKey ContextKey = "userID"

type respWriterStatusCode struct {
	http.ResponseWriter
	status int
}

func (rw *respWriterStatusCode) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// CORS Access-Control-Allow-Origin
func CORS(cc config.CORSConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			origin := cors.GetOrigin(r)
			if !cors.IsAllowed(origin, cc.Origins) {
				place := "middleware/CORS"
				utils.PrintResult(re.ErrorCORS(origin), http.StatusForbidden, place)
				rw.WriteHeader(http.StatusForbidden)
				utils.Debug(false, "cors no!!!!!!!!!!")
				return
			}

			utils.Debug(false, "cors allow!")
			cors.SetCORS(rw, cc, origin)
			next.ServeHTTP(rw, r)
			return
		})
	}
}

// Auth Check cookie exists
func Auth(cc config.SessionConfig, oauth oauth2.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var (
				maxLimit = 3
				userID   string
				err      error
			)
			for i := 0; i < maxLimit; i++ {
				userID, err = auth.Check(rw, r, oauth)
				if err == nil {
					break
				}
			}
			if err != nil {
				const place = "middleware/Auth"
				rw.WriteHeader(http.StatusUnauthorized)
				utils.PrintResult(err, http.StatusUnauthorized, place)
				utils.SendErrorJSON(rw, http.ErrNoCookie, place)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserKey, userID)

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
