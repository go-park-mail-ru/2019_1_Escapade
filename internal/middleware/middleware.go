package middleware

import (
	config "escapade/internal/config"
	cors "escapade/internal/cors"
	cookie "escapade/internal/misc"
	"fmt"
	"log"

	"net/http"
	"strings"
)

// HandleDecorator middleware
type HandleDecorator func(http.HandlerFunc) http.HandlerFunc

func setCORS(rw http.ResponseWriter, cc config.CORSConfig, name string) {
	rw.Header().Set("Access-Control-Allow-Origin", name)
	rw.Header().Set("Access-Control-Allow-Headers", strings.Join(cc.Headers, ", "))
	rw.Header().Set("Access-Control-Allow-Credentials", cc.Credentials)
	rw.Header().Set("Access-Control-Allow-Methods", strings.Join(cc.Methods, ", "))
}

// CORS Access-Control-Allow-Origin
func CORS(cc config.CORSConfig, preCORS bool) HandleDecorator {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if !cors.IsAllowed(origin, cc.Origins) {
				fmt.Println("CORS doesnt want work with you, mr " + origin)
				rw.WriteHeader(http.StatusForbidden)
				return
			}
			setCORS(rw, cc, origin)

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
func Auth() HandleDecorator {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookie.NameCookie)

			if err != nil || cookie.Value == "" {

				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			next(rw, r)
		}
	}
}

//Recover catch panic
func Recover() HandleDecorator {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic: %+v", err)
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next(rw, r)
		}
	}
}

// ApplyMiddleware apply middleware
func ApplyMiddleware(handler http.HandlerFunc,
	decorators ...HandleDecorator) http.HandlerFunc {
	handler = Recover()(handler)
	for _, m := range decorators {
		handler = m(handler)
	}
	return handler
}
