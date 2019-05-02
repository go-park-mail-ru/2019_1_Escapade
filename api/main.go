package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2019_1_Escapade/api/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"net/http"

	"go.uber.org/zap"
)

// ./swag init

// @title Escapade API
// @version 1.0
// @description Documentation

// @host https://escapade-backend.herokuapp.com
// @BasePath /api/v1
func main() {
	const (
		place      = "main"
		confPath   = "conf.json"
		secretPath = "secret.json"
	)

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Zap logger error:", err)
	}
	defer logger.Sync()

	authConn, err := clients.ServiceConnectionsInit()
	if err != nil {
		log.Fatal("serviceConnectionsInit error:", err)
	}
	defer authConn.Close()

	API, Conf, err := api.InitAPI(confPath, secretPath, authConn) // init.go

	if err != nil {
		utils.PrintResult(err, 0, "main")
		return
	}
	API.RandomUsers(10) // create 10 users for tests
	r := router.GetRouter(API, Conf)
	port := router.GetPort(Conf)

	game.Launch(&Conf.Game)
	defer game.GetLobby().Stop()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		game.GetLobby().Stop()
		game.GetLobby().Free()
		os.Exit(1)
	}()

	if err = http.ListenAndServe(port, r); err != nil {
		utils.PrintResult(err, 0, "main")
	}
}
