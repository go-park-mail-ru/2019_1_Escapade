package main

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
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
		confPath   = "history/history.json"
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
	r.HandleFunc("/history/ws", handler.GameHistory)
	r.Handle("/history/metrics", promhttp.Handler())

	metrics.InitHitsMetric("api")

	prometheus.MustRegister(metrics.Hits)

	port := router.GetPort(configuration)
	if err = http.ListenAndServe(port, r); err != nil {
		utils.PrintResult(err, 0, "main")
	}
}
