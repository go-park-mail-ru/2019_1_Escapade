package middleware

import (
	config "escapade/internal/config"
	cors "escapade/internal/cors"
	cookie "escapade/internal/misc"
	"fmt"

	"net/http"
	"strings"
)

type handleDecorator func(http.HandlerFunc) http.HandlerFunc

func setCORS(rw http.ResponseWriter, cc config.CORSConfig) {
	rw.Header().Set("Access-Control-Allow-Origin", strings.Join(cc.Origins, ", "))
	rw.Header().Set("Access-Control-Allow-Headers", strings.Join(cc.Headers, ", "))
	rw.Header().Set("Access-Control-Allow-Credentials", cc.Credentials)
	rw.Header().Set("Access-Control-Allow-Methods", strings.Join(cc.Methods, ", "))
}

// CORS Access-Control-Allow-Origin
func CORS(cc config.CORSConfig) handleDecorator {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if !cors.IsAllowed(origin, cc.Origins) {
				fmt.Println("CORS doesnt want work with you, mr " + origin)
				rw.WriteHeader(http.StatusForbidden)
				return
			}
			setCORS(rw, cc)
			fmt.Println("CORS disabled")

			hf(rw, r)
			return
		}
	}
}

// В будущем отрефакторить, ибо явно дублирует CORS
// PreflightRequest
func PRCORS(cc config.CORSConfig) handleDecorator {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if !cors.IsAllowed(origin, cc.Origins) {
				rw.WriteHeader(http.StatusForbidden)
				return
			}

			setCORS(rw, cc)

			rw.WriteHeader(http.StatusOK)

			return
		}
	}
}

// Check cookie
func Auth() handleDecorator {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookie.NameCookie)

			if err != nil || cookie.Value == "" {

				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			hf(rw, r)
		}
	}
}
