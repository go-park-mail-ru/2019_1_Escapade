package middleware

import (
	"context"
	"escapade/internal/config"
	cookie "escapade/internal/cookie"
	"escapade/internal/cors"
	re "escapade/internal/return_errors"
	"escapade/internal/utils"

	"net/http"

	"go.uber.org/zap"
)

// HandleDecorator middleware
type HandleDecorator func(http.HandlerFunc) http.HandlerFunc
type ApplyMiddleware func(handler http.HandlerFunc, decorators ...HandleDecorator) http.HandlerFunc

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
func Auth(cc config.CookieConfig) HandleDecorator {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			if _, err := cookie.GetSessionCookie(r, cc); err != nil {
				const place = "middleware/Auth"
				utils.PrintResult(err, http.StatusUnauthorized, place)
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

//Log
func Logger(next http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		newLogger := logger.With(zap.String("method", r.Method), zap.String("url", r.URL.Path))
		newLogger.Info("Request Started")
		ctx := r.Context()
		ctx = context.WithValue(ctx, "logger", newLogger)
		next(rw, r.WithContext(ctx))
		newLogger.Info("Request Ended")
	}
}

// ApplyMiddleware apply middleware
func ApplyMiddlewareLogger(logger *zap.Logger) ApplyMiddleware {
	return func(handler http.HandlerFunc,
		decorators ...HandleDecorator) http.HandlerFunc {
		handler = Recover(handler)
		for _, m := range decorators {
			handler = m(handler)
		}
		handler = Logger(handler, logger)
		return handler
	}
}
