package main

import (
	"escapade/internal/config"
	"escapade/internal/database"
	mi "escapade/internal/middleware"
	"escapade/internal/services/api"
	"fmt"
	"os"

	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	confPath = "conf.json"
)

// @title Escapade API
// @version 1.0
// @description Documentation

// @host https://escapade-backend.herokuapp.com
// @BasePath /api/v1
func main() {

	fmt.Println("Ok")
	conf, confErr := config.Init(confPath)

	fmt.Println("Ok")
	if confErr != nil {
		panic(confErr)
	}

	fmt.Println("Ok")
	db, dbErr := database.Init(conf.DataBase)
	if dbErr != nil {
		panic(dbErr)
	}

	fmt.Println("Ok")
	API := api.Init(db, conf.Storage)

	r := mux.NewRouter()

	r.PathPrefix("/api/v1/")

	r.HandleFunc("/", mi.CORS(conf.Cors)(API.Ok))
	r.HandleFunc("/user", mi.CORS(conf.Cors)(API.GetMyProfile)).Methods("GET")
	r.HandleFunc("/user", mi.CORS(conf.Cors)(API.CreateUser)).Methods("POST")
	r.HandleFunc("/user", mi.CORS(conf.Cors)(API.DeleteAccount)).Methods("DELETE")
	r.HandleFunc("/user", mi.CORS(conf.Cors)(API.UpdateProfile)).Methods("PUT")
	r.HandleFunc("/user", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	r.HandleFunc("/session", mi.CORS(conf.Cors)(API.Logout)).Methods("DELETE")
	r.HandleFunc("/session", mi.CORS(conf.Cors)(API.Login)).Methods("POST")
	r.HandleFunc("/session", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	r.HandleFunc("/avatar", mi.CORS(conf.Cors)(API.GetImage)).Methods("GET")
	r.HandleFunc("/avatar", mi.CORS(conf.Cors)(API.PostImage)).Methods("POST")
	r.HandleFunc("/avatar", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	r.HandleFunc("/users", mi.CORS(conf.Cors)(API.GetUsers)).Methods("GET")
	r.HandleFunc("/users/pages/{page}", mi.CORS(conf.Cors)(API.GetUsers)).Methods("GET")
	r.HandleFunc("/users/pages_amount", mi.CORS(conf.Cors)(API.GetUsersPageAmount)).Methods("GET")

	r.HandleFunc("/users/{name}/games", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	r.HandleFunc("/users/{name}/games/{page}", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	r.HandleFunc("/users/{name}/profile", mi.CORS(conf.Cors)(API.GetProfile)).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	fmt.Println("launched, look at us on " + conf.Server.Host + conf.Server.Port) //+ os.Getenv("PORT"))

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "3000")
	}

	err := http.ListenAndServe(":"+os.Getenv("PORT"), r)
	fmt.Println("oh, this is error:" + err.Error())
}
