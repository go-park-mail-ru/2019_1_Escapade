package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"google.golang.org/grpc"

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

	authConn := serviceConnectionsInit()
	defer authConn.Close()

	API, conf, err := api.GetHandler(confPath, secretPath) // init.go

	API.Clients = clients.Init(authConn)
	if err != nil {
		utils.PrintResult(err, 0, "main")
		return
	}
	API.RandomUsers(10) // create 10 users for tests
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

func serviceConnectionsInit() (authConn *grpc.ClientConn) {
	if os.Getenv("AUTHSERVICE_URL") == "" {
		os.Setenv("AUTHSERVICE_URL", "localhost:3333")
	}
	authConn, err := grpc.Dial(
		os.Getenv("AUTHSERVICE_URL"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Cant connect to auth service!")
	}

	//Other micro services conns wiil be here

	return
}
