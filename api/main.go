package main

import (
	"escapade/internal/config"
	"escapade/internal/router"
	"escapade/internal/services/api"
	"escapade/internal/utils"
	"log"
	"os"

	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc"
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
		confPath   = "./conf.json"
		secretPath = "./secret.json"
	)

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Zap logger error:", err)
	}
	defer logger.Sync()

	authConn := serviceConnectionsInit()
	defer authConn.Close()

	var (
		conf *config.Configuration
	)

	api.API, conf, err = api.GetHandler(confPath, secretPath, authConn) // init.go
	api.API.DB.RandomUsers(10)                                              // create 10 users for tests
	if err != nil {
		utils.PrintResult(err, 0, "main")
		return
	}
	r := router.GetRouter(api.API, conf, logger)
	port := router.GetPort(conf)

	//server := grpc.NewServer()
	//session.RegisterAuthCheckerServer(server, sessMan.NewSessionManager(redisConn))

	logger.Info("Starting server",
		zap.String("port", port))

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
