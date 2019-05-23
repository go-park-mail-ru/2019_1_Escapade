package router

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"fmt"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// GetRouter return router
func GetRouter(API *api.Handler, conf *config.Configuration) *mux.Router {
	r := mux.NewRouter()

	var v = r.PathPrefix("/api").Subrouter()

	v.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	var v1 = v

	v1.HandleFunc("/", mi.ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, false))).Methods("GET")
	//r.HandleFunc("/ws", mi.ApplyMiddleware(API.GameOnline,
	//	mi.CORS(conf.Cors, false)))

	v1.HandleFunc("/user", mi.ApplyMiddleware(API.GetMyProfile,
		mi.Auth(conf.Session), mi.CORS(conf.Cors, false))).Methods("GET")
	v1.HandleFunc("/user", mi.ApplyMiddleware(API.CreateUser,
		mi.CORS(conf.Cors, false))).Methods("POST")
	v1.HandleFunc("/user", mi.ApplyMiddleware(API.DeleteUser,
		mi.Auth(conf.Session), mi.CORS(conf.Cors, false))).Methods("DELETE")
	v1.HandleFunc("/user", mi.ApplyMiddleware(API.UpdateProfile,
		mi.Auth(conf.Session), mi.CORS(conf.Cors, false))).Methods("PUT")
	v1.HandleFunc("/user", mi.ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	v1.HandleFunc("/session", mi.ApplyMiddleware(API.Logout,
		mi.CORS(conf.Cors, false))).Methods("DELETE")
	v1.HandleFunc("/session", mi.ApplyMiddleware(API.Login,
		mi.CORS(conf.Cors, false))).Methods("POST")
	v1.HandleFunc("/session", mi.ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	v1.HandleFunc("/avatar/{name}", mi.ApplyMiddleware(API.GetImage,
		mi.CORS(conf.Cors, false))).Methods("GET")
	v1.HandleFunc("/avatar/{name}", mi.ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	v1.HandleFunc("/avatar", mi.ApplyMiddleware(API.PostImage,
		mi.Auth(conf.Session), mi.CORS(conf.Cors, false))).Methods("POST")
	v1.HandleFunc("/avatar", mi.ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	v1.HandleFunc("/users/pages", mi.ApplyMiddleware(API.GetUsers,
		mi.CORS(conf.Cors, false))).Methods("GET")
	v1.HandleFunc("/users/pages", mi.ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")
	v1.HandleFunc("/users/pages_amount", mi.ApplyMiddleware(API.GetUsersPageAmount,
		mi.CORS(conf.Cors, false))).Methods("GET")

	v1.HandleFunc("/game", mi.ApplyMiddleware(API.SaveRecords,
		mi.Auth(conf.Session), mi.CORS(conf.Cors, false))).Methods("POST")
	v1.HandleFunc("/game", mi.ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	// v1.HandleFunc("/users/{name}/games", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	// v1.HandleFunc("/users/{name}/games/{page}", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	// v1.HandleFunc("/users/{name}/profile", mi.CORS(conf.Cors)(API.GetProfile)).Methods("GET")

	r.Handle("/metrics", promhttp.Handler())
	return r
}

// GetPort return port
func GetPort(conf *config.Configuration) (port string) {

	env := os.Getenv(conf.Server.PortURL)
	if os.Getenv(conf.Server.PortURL)[0] != ':' {
		port = ":" + env
	} else {
		port = env
	}
	fmt.Println("launched, look at us on " + conf.Server.Host + port)
	return
}
