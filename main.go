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
)

const (
	confPath = "conf.json"
)

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

	r.HandleFunc("/", mi.CORS(conf.Cors)(API.Ok))
	r.HandleFunc("/register", mi.CORS(conf.Cors)(API.Register)).Methods("POST")
	r.HandleFunc("/register", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	r.HandleFunc("/delete", mi.CORS(conf.Cors)(API.DeleteAccount)).Methods("DELETE")
	r.HandleFunc("/delete", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	r.HandleFunc("/logout", mi.CORS(conf.Cors)(API.Logout)).Methods("DELETE")
	r.HandleFunc("/logout", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	r.HandleFunc("/login", mi.CORS(conf.Cors)(API.Login)).Methods("POST")
	r.HandleFunc("/login", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	r.HandleFunc("/playerAvatar", mi.CORS(conf.Cors)(API.GetImage)).Methods("GET")
	r.HandleFunc("/playerAvatar", mi.CORS(conf.Cors)(API.PostImage)).Methods("POST")
	r.HandleFunc("/playerAvatar", mi.PRCORS(conf.Cors)(API.Ok)).Methods("OPTIONS")

	r.HandleFunc("/users", mi.CORS(conf.Cors)(API.GetUsers)).Methods("GET")
	r.HandleFunc("/users/{page}", mi.CORS(conf.Cors)(API.GetUsers)).Methods("GET")
	r.HandleFunc("/usersPageAmount", mi.CORS(conf.Cors)(API.GetUsersPageAmount)).Methods("GET")

	r.HandleFunc("/me", mi.CORS(conf.Cors)(API.Me)).Methods("GET")
	r.HandleFunc("/{name}/games", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	r.HandleFunc("/{name}/games/{page}", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	r.HandleFunc("/{name}", mi.CORS(conf.Cors)(API.GetProfile)).Methods("GET")

	fmt.Println("launched, look at us on " + conf.Server.Host + conf.Server.Port) //+ os.Getenv("PORT"))

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "3000")
	}

	err := http.ListenAndServe(":"+os.Getenv("PORT"), r)
	fmt.Println("oh, this is error:" + err.Error())
}
