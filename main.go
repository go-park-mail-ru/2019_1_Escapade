package main

import (
	"escapade/internal/router"
	"escapade/internal/services/api"
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
	confPath := router.GetConf() //"test.json"

	fmt.Println("we use configuration", confPath)

	API, conf, err := api.GetHandler(confPath) // init.go
	if err != nil {
		fmt.Println("Some error happened with configuration file or database" + err.Error())
		return
	}
	r := router.GetRouter(API, conf)
	port := router.GetPort(conf)

	if err = http.ListenAndServe(port, r); err != nil {
		fmt.Println("Server error:" + err.Error())
	}
}
