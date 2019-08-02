package main

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/game"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"os"
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
			len(os.Args)-1, "get.")
		return
	}

	var (
		configurationPath = os.Args[1]
		photoPublicPath   = os.Args[2]
		photoPrivatePath  = os.Args[3]
		fieldPath         = os.Args[4]
		roomPath          = os.Args[5]
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

	metrics.InitApi()
	metrics.InitGame()

	readyChan := make(chan error)
	finishChan := make(chan interface{})

	clients.ALL.Init(readyChan, finishChan, configuration.Services...)

	handler = api.Init(db, configuration)

	game.Launch(&configuration.Game, db, photo.GetImages)

	var (
		r    = server.GameRouter(handler, configuration.Cors)
		port = server.Port(configuration)
		srv  = server.Server(r, configuration.Server, port)
	)

	server.LaunchHTTP(srv, configuration.Server, func() {
		finishChan <- nil
		game.GetLobby().Stop()
		db.Db.Close()
	})
	os.Exit(0)
}
