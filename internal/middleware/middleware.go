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

type handleDecorator func(http.HandlerFunc) http.HandlerFunc

func setCORS(rw http.ResponseWriter, cc config.CORSConfig) {
	rw.Header().Set("Access-Control-Allow-Origin", strings.Join(cc.Origins, ", "))
	rw.Header().Set("Access-Control-Allow-Headers", strings.Join(cc.Headers, ", "))
	rw.Header().Set("Access-Control-Allow-Credentials", cc.Credentials)
	rw.Header().Set("Access-Control-Allow-Methods", strings.Join(cc.Methods, ", "))
}

// CORS Access-Control-Allow-Origin
func CORS(cc config.CORSConfig, preCORS bool) handleDecorator {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if !cors.IsAllowed(origin, cc.Origins) {
				fmt.Println("CORS doesnt want work with you, mr " + origin)
				rw.WriteHeader(http.StatusForbidden)
				return
			}
			setCORS(rw, cc)
			fmt.Println("CORS disabled")

			if preCORS {
				rw.WriteHeader(http.StatusOK)
			} else {
				next(rw, r)
			}
			return
		}
	}
}

// В будущем отрефакторить, ибо явно дублирует CORS
// PreflightRequest
// func PRCORS(cc config.CORSConfig) handleDecorator {
// 	return func(next http.HandlerFunc) http.HandlerFunc {
// 		return func(rw http.ResponseWriter, r *http.Request) {
// 			origin := r.Header.Get("Origin")

// 			if !cors.IsAllowed(origin, cc.Origins) {
// 				rw.WriteHeader(http.StatusForbidden)
// 				return
// 			}

// 			setCORS(rw, cc)

// 			rw.WriteHeader(http.StatusOK)

// 			return
// 		}
// 	}
// }

// Check cookie
func Auth() handleDecorator {
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

func Recover() handleDecorator {
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

// ApplyDecorators
func ApplyMiddleware(handler http.HandlerFunc,
	decorators ...handleDecorator) http.HandlerFunc {
	handler = Recover()(handler)
	for _, m := range decorators {
		handler = m(handler)
	}
	return handler
}
