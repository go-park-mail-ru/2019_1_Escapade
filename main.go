package main

import (
	"escapade/internal/router"
	"escapade/internal/services/api"
	"escapade/internal/services/game"
	"escapade/internal/utils"
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
		place    = "main"
		confPath = "conf.json"
	)

	API, conf, err := api.GetHandler(confPath) // init.go
	API.DB.RandomUsers(10)                     // create 10 users for tests
	if err != nil {
		utils.PrintResult(err, 0, "main")
		return
	}
	r := router.GetRouter(API, conf)
	port := router.GetPort(conf)

	game.Launch(&conf.Game)
	defer game.GetLobby().Stop()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		game.GetLobby().Free()
		os.Exit(1)
	}()

	if err = http.ListenAndServe(port, r); err != nil {
		utils.PrintResult(err, 0, "main")
	}
}
