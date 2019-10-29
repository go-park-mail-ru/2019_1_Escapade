package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	server "github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	//
	_ "net/http/pprof"

	_ "github.com/go-park-mail-ru/2019_1_Escapade/docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router return router of api operations
func Router(API *Handler, cors config.CORS, cc config.Cookie,
	ca config.Auth, client config.AuthClient) *mux.Router {

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
	api.HandleFunc("/users/pages_amount", API.HandleUsersPageAmount).Methods("GET")

	r.PathPrefix("/health").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		v, _ := server.GetIP()
		fmt.Println("fun:", v)
		rw.Write([]byte("all ok " + v.String()))
	})

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	router.Use(r, mi.CORS(cors))
	apiWithAuth.Use(mi.Auth(cc, ca, client))
	return r
}

// HistoryRouter return router for history service
/*
func HistoryRouter(handler *api.Handler, cors config.CORS) *mux.Router {
	r := mux.NewRouter()

	var history = r.PathPrefix("/history").Subrouter()

	history.Use(mi.Recover, mi.CORS(cors), mux.CORSMethodMiddleware(r))

	history.HandleFunc("/ws", handler.GameOnline)
	history.Handle("/metrics", promhttp.Handler())
	return r
}
*/
