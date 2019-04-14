package main

import (
	mi "escapade/internal/middleware"
	"escapade/internal/services/api"
	"fmt"
	"os"

	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// ./swag init

// @title Escapade API
// @version 1.0
// @description Documentation

// @host https://escapade-backend.herokuapp.com
// @BasePath /api/v1
func main() {
	const confPath = "conf.json"

	API, conf, err := api.GetHandler(confPath) // init.go
	if err != nil {
		fmt.Println("Some error happened with configuration file or database" + err.Error())
		return
	}

	r := mux.NewRouter()

	var v = r.PathPrefix("/api").Subrouter()

	v.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	var v1 = r

	v1.HandleFunc("/", mi.CORS(conf.Cors)(API.Ok))
	r.HandleFunc("/ws", mi.CORS(conf.Cors)(API.GameOnline))

	v1.HandleFunc("/user", mi.CORS(conf.Cors)(API.GetMyProfile)).Methods("GET")
	v1.HandleFunc("/user", mi.CORS(conf.Cors)(API.CreateUser)).Methods("POST")
	v1.HandleFunc("/user", mi.CORS(conf.Cors)(API.DeleteUser)).Methods("DELETE")
	v1.HandleFunc("/user", mi.CORS(conf.Cors)(API.UpdateProfile)).Methods("PUT")
	v1.HandleFunc("/user", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	v1.HandleFunc("/session", mi.CORS(conf.Cors)(API.Logout)).Methods("DELETE")
	v1.HandleFunc("/session", mi.CORS(conf.Cors)(API.Login)).Methods("POST")
	v1.HandleFunc("/session", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	v1.HandleFunc("/avatar/{name}", mi.CORS(conf.Cors)(API.GetImage)).Methods("GET")
	v1.HandleFunc("/avatar/{name}", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	v1.HandleFunc("/avatar", mi.CORS(conf.Cors)(API.PostImage)).Methods("POST")
	v1.HandleFunc("/avatar", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	v1.HandleFunc("/users", mi.CORS(conf.Cors)(API.GetUsers)).Methods("GET")
	v1.HandleFunc("/users/pages", mi.CORS(conf.Cors)(API.GetUsers)).Methods("GET")
	v1.HandleFunc("/users/pages_amount", mi.CORS(conf.Cors)(API.GetUsersPageAmount)).Methods("GET")
	v1.HandleFunc("/avatar/pages", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")
	v1.HandleFunc("/avatar/pages_amount", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	v1.HandleFunc("/users/{name}/games", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	v1.HandleFunc("/users/{name}/games/{page}", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	v1.HandleFunc("/users/{name}/profile", mi.CORS(conf.Cors)(API.GetProfile)).Methods("GET")

	fmt.Println("launched, look at us on " + conf.Server.Host + conf.Server.Port) //+ os.Getenv("PORT"))

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", conf.Server.Port)
	}

	fmt.Println("launched, look at us on " + conf.Server.Host + os.Getenv("PORT")) //+ os.Getenv("PORT"))
	var port string
	if os.Getenv("PORT")[0] != ':' {
		port = ":" + os.Getenv("PORT")
	} else {
		port = os.Getenv("PORT")
	}

	if err = http.ListenAndServe(port, r); err != nil {
		fmt.Println("oh, this is error:" + err.Error())
	}
}
