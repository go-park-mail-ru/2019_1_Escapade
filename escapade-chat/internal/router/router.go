package router

import (
	"escapade/internal/config"
	mi "escapade/internal/middleware"
	api "escapade/internal/services/game"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// GetRouter return router
func GetRouter(API *api.Handler, conf *config.Configuration, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	var v = r.PathPrefix("/api").Subrouter()

	v.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	var v1 = r
	ApplyMiddleware := mi.ApplyMiddlewareLogger(logger)
	v1.HandleFunc("/", ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, false))).Methods("GET")
	r.HandleFunc("/ws", API.Chat)

	return r
}

// GetPort return port
func GetPort(conf *config.Configuration) (port string) {
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", conf.Server.Port)
	}

	if os.Getenv("PORT")[0] != ':' {
		port = ":" + os.Getenv("PORT")
	} else {
		port = os.Getenv("PORT")
	}
	return
}
