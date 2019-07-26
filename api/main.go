package main

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	_ "github.com/go-park-mail-ru/2019_1_Escapade/docs"

	"net/http"
	"os"
)

// to generate docs, call from root "swag init -g api/main.go"

// @title Escapade Explosion API
// @version 1.0
// @description API documentation

// @host https://explosion.team
// @BasePath /api
func main() {
	var (
		configuration *config.Configuration
		API           *api.Handler
		err           error
	)

	if len(os.Args) < 4 {
		utils.Debug(false, "Api service need 3 command line arguments. But",
			len(os.Args)-1, "get.")
		return
	}

	var (
		configurationPath = os.Args[1]
		photoPublicPath   = os.Args[2]
		photoPrivatePath  = os.Args[3]
	)

	configuration, err = config.Init(configurationPath)
	if err != nil {
		utils.Debug(false, "Initialization error with main configuration:", err.Error())
		return
	}

	err = photo.Init(photoPublicPath, photoPrivatePath)
	if err != nil {
		utils.Debug(false, "Initialization error with photo configuration:", err.Error())
		return
	}

	API, err = api.GetHandler(configuration)
	if err != nil {
		utils.Debug(false, "Initialization error with photo configuration:", err.Error())
		return
	}

	metrics.InitApi()

	/*
		authConn, err := clients.ServiceConnectionsInit(configuration.AuthClient)
		if err != nil {
			log.Fatal("serviceConnectionsInit error:", err)
		}
		defer authConn.Close()
	*/

	//API.RandomUsers(10) // create 10 users for tests
	var (
		r    = router.APIRouter(API, configuration.Cors, configuration.Session)
		port = router.Port(configuration)
	)

	if err = http.ListenAndServe(port, r); err != nil {
		utils.PrintResult(err, 0, "main")
	}
}
