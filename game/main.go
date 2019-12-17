package main

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/game"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"

	// mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// conf := config.Configure(confPath)
	// curlog := lg.Construct(logPath, logFile)
	// db := database.New(curlog)

	const (
		place      = "main"
		confPath   = "game/game.json"
		secretPath = "secret.json"
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
	config.InitPrivate(secretPath)

	metrics.InitRoomMetric("game")
	metrics.InitPlayersMetric("game")

	prometheus.MustRegister(metrics.Rooms, metrics.Players, metrics.FreeRooms, metrics.WaitingPlayers)
	/*
		authConn, err := clients.ServiceConnectionsInit(configuration.AuthClient)
		if err != nil {
			log.Fatal("serviceConnectionsInit error:", err)
		}
		defer authConn.Close()
	*/

	handler, err = api.GetGameHandler(configuration /*, authConn*/) // init.go
	if err != nil {
		fmt.Println("eeeer", err.Error())
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/game/ws", handler.GameOnline)

	r.Handle("/game/metrics", promhttp.Handler())

	game.Launch(&configuration.Game, &handler.DB, handler.Setfiles)
	defer game.GetLobby().Stop()

	// c := make(chan os.Signal)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-c
	// 	game.GetLobby().Stop()
	// 	game.GetLobby().Free()
	// 	os.Exit(1)
	// }()

	port := router.GetPort(configuration)
	if err = http.ListenAndServe(port, r); err != nil {
		utils.PrintResult(err, 0, "main")
	}
}
