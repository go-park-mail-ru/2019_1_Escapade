package main

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game"
	"os"
	"os/signal"
	"syscall"

	"net/http"
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

	API, conf, err := api.GetHandler(confPath, secretPath) // init.go

	if err != nil {
		utils.PrintResult(err, 0, "main")
		return
	}
	API.DB.RandomUsers(10) // create 10 users for tests
	r := router.GetRouter(API, conf)
	port := router.GetPort(conf)

	game.Launch(&conf.Game)
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
