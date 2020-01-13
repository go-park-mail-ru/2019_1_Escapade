package middleware

import (
	"context"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
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
			origin := GetOrigin(r)
			utils.Debug(false, "cors check")
			if !IsAllowed(origin, cc.Origins) {
				handlers.SendResult(rw, handlers.NewResult(http.StatusForbidden, "cors", nil, re.CORS(origin)))
				utils.Debug(false, "cors no!!!!!!!!!!")
				return
			}
			utils.Debug(false, "cors allow!")
			SetCORS(rw, cc, origin)
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
					utils.Debug(false, "no error auth")
					break
				} else {
					utils.Debug(false, "error auth", err.Error())
				}
			}
			if err != nil {
				const place = "middleware/Auth"
				handlers.SendResult(rw, handlers.NewResult(http.StatusUnauthorized, place, nil, err))
				return
			}
			ctx := context.WithValue(r.Context(), handlers.ContextUserKey, userID)

			utils.Debug(false, "auth end", userID, handlers.ContextUserKey)
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

//Logger log request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ip := server.GetIP(nil)
		utils.Debug(false, "listen for you on "+ip)
		next.ServeHTTP(rw, r)
	})
}

// Metrics record metrics
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		goodRW := &respWriterStatusCode{rw, 200}
		next.ServeHTTP(goodRW, r)
		utils.Debug(false, "metrics get "+utils.String(goodRW.status))
		Hits.WithLabelValues(utils.String(goodRW.status), r.URL.Path, r.Method).Inc()
	})
}
