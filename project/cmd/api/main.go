package main

import (
	"escapade/internal/config"
	"escapade/internal/database"
	mi "escapade/internal/middleware"
	"escapade/internal/services/api"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
)

const (
	confPath = "/home/artyom/projest/back/2019_1_Escapade/conf.json"
)

func main() {

	conf, confErr := config.Init(confPath)

	if confErr != nil {
		panic(confErr)
	}

	db, dbErr := database.Init()
	if dbErr != nil {
		panic(dbErr)
	}

	API := api.Init(db)

	r := mux.NewRouter()

	r.HandleFunc("/", mi.CORS(conf.Cors)(API.Ok))
	r.HandleFunc("/register", mi.CORS(conf.Cors)(API.Register)).Methods("POST")
	r.HandleFunc("/delete", mi.CORS(conf.Cors)(API.DeleteAccount)).Methods("DELETE")
	r.HandleFunc("/delete", mi.CORS(conf.Cors)(API.DeleteAccountOptions)).Methods("OPTIONS")

	r.HandleFunc("/login", mi.CORS(conf.Cors)(API.Login)).Methods("POST")
	r.HandleFunc("/{name}/games", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")

	fmt.Println("launched, look at us on " + conf.Server.Host + conf.Server.Port)
	http.ListenAndServe(conf.Server.Port, r)
}
