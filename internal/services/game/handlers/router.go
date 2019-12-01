package handlers

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/router"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/middleware"
)

// Router return router of game operations
func (h *GameHandler) Router() *mux.Router {
	r := mux.NewRouter()

	var game = r.PathPrefix("/game").Subrouter()

	game.HandleFunc("/ws", h.Handle)
	game.Handle("/metrics", promhttp.Handler())

	router.Use(r, mi.CORS(h.c.Cors))
	return r
}
