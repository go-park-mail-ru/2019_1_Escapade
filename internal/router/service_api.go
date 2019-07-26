package router

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"fmt"
	"os"

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

	api.Use(mi.Recover, mi.CORS(cors), mi.Metrics)
	apiWithAuth.Use(mi.Recover, mi.CORS(cors), mi.Auth(session), mi.Metrics)

	api.HandleFunc("/user", API.HandleUser).Methods("POST", "GET", "OPTIONS")
	apiWithAuth.HandleFunc("/user", API.HandleUser).Methods("DELETE", "PUT")

	api.HandleFunc("/session", API.HandleSession).Methods("POST", "OPTIONS")
	apiWithAuth.HandleFunc("/session", API.HandleSession).Methods("DELETE")

	api.HandleFunc("/avatar/{name}", API.HandleAvatar).Methods("GET")
	api.HandleFunc("/avatar", API.HandleAvatar).Methods("OPTIONS")
	apiWithAuth.HandleFunc("/avatar", API.HandleAvatar).Methods("POST")

	api.HandleFunc("/game", API.HandleOfflineGame).Methods("OPTIONS")
	apiWithAuth.HandleFunc("/game", API.HandleOfflineGame).Methods("POST")

	api.HandleFunc("/users/pages", API.HandleUsersPages).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/pages_amount", API.GetUsersPageAmount).Methods("GET")

	return r
}

// Port return port
func Port(conf *config.Configuration) (port string) {

	env := os.Getenv(conf.Server.PortURL)
	if os.Getenv(conf.Server.PortURL)[0] != ':' {
		port = ":" + env
	} else {
		port = env
	}
	fmt.Println("launched, look at us on " + conf.Server.Host + port)
	return
}
