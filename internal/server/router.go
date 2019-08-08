package server

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/go-park-mail-ru/2019_1_Escapade/docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// APIRouter return router for api
func APIRouter(API *api.Handler, cors config.CORSConfig, session config.SessionConfig) *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	var api = r.PathPrefix("/api").Subrouter()
	var apiWithAuth = r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/user", API.HandleUser).Methods("OPTIONS", "POST")
	apiWithAuth.HandleFunc("/user", API.HandleUser).Methods("DELETE", "PUT", "GET")

	api.HandleFunc("/session", API.HandleSession).Methods("POST", "OPTIONS", "DELETE")

	api.HandleFunc("/avatar/{name}", API.HandleAvatar).Methods("GET")
	api.HandleFunc("/avatar", API.HandleAvatar).Methods("OPTIONS")
	apiWithAuth.HandleFunc("/avatar", API.HandleAvatar).Methods("POST")

	api.HandleFunc("/game", API.HandleOfflineGame).Methods("OPTIONS")
	apiWithAuth.HandleFunc("/game", API.HandleOfflineGame).Methods("POST")

	api.HandleFunc("/users/pages", API.HandleUsersPages).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/pages_amount", API.GetUsersPageAmount).Methods("GET")

	r.Use(mi.Recover, mi.CORS(cors), mi.Metrics)
	apiWithAuth.Use(mi.Auth(session, API.Oauth))

	return r
}

// GameRouter return router for game
func GameRouter(handler *api.Handler, cors config.CORSConfig) *mux.Router {
	r := mux.NewRouter()

	var game = r.PathPrefix("/game").Subrouter()

	game.Use(mi.Recover, mi.CORS(cors))

	game.HandleFunc("/ws", handler.GameOnline)
	game.Handle("/metrics", promhttp.Handler())
	return r
}

// HistoryRouter return router for history service
func HistoryRouter(handler *api.Handler, cors config.CORSConfig) *mux.Router {
	r := mux.NewRouter()

	var history = r.PathPrefix("/history").Subrouter()

	history.Use(mi.Recover, mi.CORS(cors), mux.CORSMethodMiddleware(r))

	history.HandleFunc("/ws", handler.GameOnline)
	history.Handle("/metrics", promhttp.Handler())
	return r
}
