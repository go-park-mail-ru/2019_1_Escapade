package main

import (
	"escapade/internal/router"
	"escapade/internal/services/api"
	"escapade/internal/utils"
	"fmt"

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
		confPath   = "./conf.json"
		secretPath = "./secret.json"
	)

	API, conf, err := api.GetHandler(confPath, secretPath) // init.go
	API.DB.RandomUsers(10)                                 // create 10 users for tests
	if err != nil {
		utils.PrintResult(err, 0, "main")
		return
	}
	r := router.GetRouter(API, conf)
	port := router.GetPort(conf)

	fmt.Println("API RUN")
	if err = http.ListenAndServe(port, r); err != nil {
		utils.PrintResult(err, 0, "main")
	}
}
