package router

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
)

// GameRouter return router for game
func GameRouter(handler *api.Handler, cors config.CORSConfig) *mux.Router {
	r := mux.NewRouter()

	var game = r.PathPrefix("/game").Subrouter()

	game.Use(mi.Recover, mi.CORS(cors), mi.Metrics)

	game.HandleFunc("/ws", handler.GameOnline)
	game.Handle("/metrics", promhttp.Handler())
	return r
}
