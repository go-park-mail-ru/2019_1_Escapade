package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
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
		confPath   = "api/api.json"
		secretPath = "secret.json"
	)

	var (
		configuration *config.Configuration
		API           *api.Handler
	)

	metrics.InitHitsMetric("api")

	prometheus.MustRegister(metrics.Hits)

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Zap logger error:", err)
		return
	}
	defer logger.Sync()

	if configuration, err = config.InitPublic(confPath); err != nil {
		fmt.Println("eeeer", err.Error())
		return
	}
	config.InitPrivate(secretPath)

	/*
		authConn, err := clients.ServiceConnectionsInit(configuration.AuthClient)
		if err != nil {
			log.Fatal("serviceConnectionsInit error:", err)
		}
		defer authConn.Close()
	*/
	API, err = api.GetAPIHandler(configuration /*, authConn*/) // init.go

	if err != nil {
		utils.PrintResult(err, 0, "main")
		return
	}
	//API.RandomUsers(10) // create 10 users for tests
	r := router.GetRouter(API, configuration)
	port := router.GetPort(configuration)

	if err = http.ListenAndServe(port, r); err != nil {
		utils.PrintResult(err, 0, "main")
	}
}
