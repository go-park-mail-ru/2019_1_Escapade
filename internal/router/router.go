package router

import (
	"escapade/internal/config"
	mi "escapade/internal/middleware"
	"escapade/internal/services/api"
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
	r.HandleFunc("/ws", ApplyMiddleware(API.GameOnline,
		mi.CORS(conf.Cors, false)))

	v1.HandleFunc("/user", ApplyMiddleware(API.GetMyProfile,
		mi.Auth(conf.Cookie), mi.CORS(conf.Cors, false))).Methods("GET")
	v1.HandleFunc("/user", ApplyMiddleware(API.CreateUser,
		mi.CORS(conf.Cors, false))).Methods("POST")
	v1.HandleFunc("/user", ApplyMiddleware(API.DeleteUser,
		mi.Auth(conf.Cookie), mi.CORS(conf.Cors, false))).Methods("DELETE")
	v1.HandleFunc("/user", ApplyMiddleware(API.UpdateProfile,
		mi.Auth(conf.Cookie), mi.CORS(conf.Cors, false))).Methods("PUT")
	v1.HandleFunc("/user", ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	v1.HandleFunc("/session", ApplyMiddleware(API.Logout,
		mi.CORS(conf.Cors, false))).Methods("DELETE")
	v1.HandleFunc("/session", ApplyMiddleware(API.Login,
		mi.CORS(conf.Cors, false))).Methods("POST")
	v1.HandleFunc("/session", ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	v1.HandleFunc("/avatar/{name}", ApplyMiddleware(API.GetImage,
		mi.CORS(conf.Cors, false))).Methods("GET")
	v1.HandleFunc("/avatar/{name}", ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	v1.HandleFunc("/avatar", ApplyMiddleware(API.PostImage,
		mi.Auth(conf.Cookie), mi.CORS(conf.Cors, false))).Methods("POST")
	v1.HandleFunc("/avatar", ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	v1.HandleFunc("/users/pages", ApplyMiddleware(API.GetUsers,
		mi.CORS(conf.Cors, false))).Methods("GET")
	v1.HandleFunc("/users/pages", ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")
	v1.HandleFunc("/users/pages_amount", ApplyMiddleware(API.GetUsersPageAmount,
		mi.CORS(conf.Cors, false))).Methods("GET")

	v1.HandleFunc("/game", ApplyMiddleware(API.SaveRecords,
		mi.Auth(conf.Cookie), mi.CORS(conf.Cors, false))).Methods("POST")
	v1.HandleFunc("/game", ApplyMiddleware(API.Ok,
		mi.CORS(conf.Cors, true))).Methods("OPTIONS")

	// v1.HandleFunc("/users/{name}/games", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	// v1.HandleFunc("/users/{name}/games/{page}", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	// v1.HandleFunc("/users/{name}/profile", mi.CORS(conf.Cors)(API.GetProfile)).Methods("GET")

	return r
}

func GetConf() string {
	if os.Getenv("PORT") == "" {
		return "conf.json"
	}
	return "deploy.json"
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
