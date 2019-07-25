package main

import (
	api "github.com/go-park-mail-ru/2019_1_Escapade/api/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/game"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"net/http"

	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		configuration *config.Configuration
		handler       *api.Handler
		db            *database.DataBase
		err           error
	)

	if len(os.Args) < 6 {
		utils.Debug(false, "Game service need 5 command line arguments. But only",
		len(os.Args), "get.")
		return
	}

	var (
		configurationPath = os.Args[1]
		photoPublicPath   = os.Args[2]
		photoPrivatePath  = os.Args[3]
		fieldPath         = os.Args[4]
		roomPath          = os.Args[5]
	)

	if configuration, err = config.Init(configurationPath); err != nil {
		utils.Debug(false, "Configuration error:", err.Error())
		return
	}
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
	
	err = constants.InitField(fieldPath)
	if err != nil {
		utils.Debug(false, "Initialization error with field constants:", err.Error())
		return
	}
	
	err = constants.InitRoom(roomPath)
	if err != nil {
		utils.Debug(false, "Initialization error with room constants:", err.Error())
		return
	}
	
	db, err = database.Init(configuration.DataBase)
	if err != nil {
		utils.Debug(false, "Initialization error with database:", err.Error())
		return
	}
	
	metrics.InitGame()
	/*
		authConn, err := clients.ServiceConnectionsInit(configuration.AuthClient)
		if err != nil {
			log.Fatal("serviceConnectionsInit error:", err)
		}
		defer authConn.Close()
	*/

	handler = api.Init(db, configuration)

	r := mux.NewRouter()
	r.HandleFunc("/game/ws", handler.GameOnline)
	r.Handle("/game/metrics", promhttp.Handler())

	game.Launch(&configuration.Game, db, photo.GetImages)
	defer game.GetLobby().Stop()

	port := router.GetPort(configuration)
	if err = http.ListenAndServe(port, r); err != nil {
		utils.Debug(false, "Serving error:", err.Error())
	}
}
