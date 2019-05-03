package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"github.com/gorilla/mux"
)

func main() {
	// conf := config.Configure(confPath)
	// curlog := lg.Construct(logPath, logFile)
	// db := database.New(curlog)

	const (
		place    = "main"
		confPath = "history/history.json"
	)

	var (
		configuration *config.Configuration
		handler       *api.Handler
		err           error
	)

	if configuration, err = config.InitPublic(confPath); err != nil {
		fmt.Println("eeeer", err.Error())
		return
	}

	authConn, err := clients.ServiceConnectionsInit(configuration.AuthClient)
	if err != nil {
		log.Fatal("serviceConnectionsInit error:", err)
	}
	defer authConn.Close()

	handler, err = api.GetGameHandler(configuration, authConn) // init.go
	if err != nil {
		fmt.Println("eeeer", err.Error())
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws", mi.ApplyMiddleware(handler.GameHistory,
		mi.CORS(configuration.Cors, false)))

	port := router.GetPort(configuration)
	if err = http.ListenAndServe(port, r); err != nil {
		utils.PrintResult(err, 0, "main")
	}
}
