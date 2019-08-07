package middleware

import (
	"context"
	"encoding/json"

	"gopkg.in/oauth2.v3/models"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/cookie"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/cors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"strconv"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
func Auth(cc config.SessionConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			/*
							timeAlive := int(token.Expiry.Sub(time.Now()).Seconds())
				http.SetCookie(rw, cookie.Cookie("accessToken", token.AccessToken, timeAlive))
				http.SetCookie(rw, cookie.Cookie("tokenType", token.TokenType, timeAlive))
				http.SetCookie(rw, cookie.Cookie("refreshToken", token.RefreshToken, timeAlive))

			*/

			cookie, err := r.Cookie("access_token")
			if err != nil || cookie == nil || cookie.Value == "" {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				utils.Debug(false, "1)something went wrong", err.Error())
				return
			}

			var authServerURL = "http://localhost:9096"
			resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s", authServerURL, cookie.Value))
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				utils.Debug(false, "2)something went wrong", err.Error())
				return
			}
			defer resp.Body.Close()

			token := models.Token{}

			err = json.NewDecoder(resp.Body).Decode(&token)

			r.Context().Value("UserID")

			if err != nil {
				const place = "middleware/Auth"
				rw.WriteHeader(http.StatusUnauthorized)
				utils.PrintResult(err, http.StatusUnauthorized, place)
				utils.SendErrorJSON(rw, re.ErrorNoCookie(), place)
				return
			} else {
				utils.Debug(false, "!!!!!!!!!!!!!!!!!!!!!token:", token.GetClientID(), token.GetCode(), token.GetScope(), token.GetAccess())
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, token.GetUserID())

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
