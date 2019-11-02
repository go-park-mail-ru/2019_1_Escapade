package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"

	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Router return router of game operations
func Router(db *database.DataBase, c *config.Configuration) *mux.Router {
	r := mux.NewRouter()

	var game = r.PathPrefix("/game").Subrouter()

	var upgraderWS = websocket.Upgrader{
		ReadBufferSize:  c.WebSocket.ReadBufferSize,
		WriteBufferSize: c.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	game.HandleFunc("/ws", gameOnline(db, c, upgraderWS))
	game.Handle("/metrics", promhttp.Handler())

	router.Use(r, mi.CORS(c.Cors))
	return r
}
