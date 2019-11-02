package main

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	ametrics "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/handlers"
	gmetrics "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"os"
)

func main() {
	var (
		configuration *config.Configuration
		db            *database.DataBase
		err           error
	)

	utils.Debug(false, "1. Check command line arguments")

	if len(os.Args) < 6 {
		utils.Debug(false, "ERROR. Game service need 5 command line arguments. But only",
			len(os.Args)-1, "get.")
		return
	}

	var (
		configurationPath = os.Args[1]
		photoPublicPath   = os.Args[2]
		photoPrivatePath  = os.Args[3]
		fieldPath         = os.Args[4]
		roomPath          = os.Args[5]
		mainPort          = os.Args[6]
		consulPort        = os.Args[7]
		mainPortInt       int
	)

	mainPort, mainPortInt, err = server.Port(mainPort)
	if err != nil {
		utils.Debug(false, "ERROR - invalid server port(cant convert to int):", err.Error())
		return
	}
	consulPort = server.FixPort(consulPort)

	utils.Debug(false, "✔")
	utils.Debug(false, "2. Set the configuration")

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

	ametrics.Init()
	gmetrics.Init()

	utils.Debug(false, "✔✔")
	utils.Debug(false, "3. Register in consul")

	var (
		serviceName = "game"
		ttl         = time.Second * 10
		maxConn     = 40
	)

	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = configuration.Server.Host
	}

	newTags := []string{"game", "traefik.frontend.entryPoints=http",
		"traefik.frontend.rule=Host:game.consul.localhost"}

	consul, err := server.InitConsulService(serviceName,
		mainPortInt, newTags, ttl, maxConn, consulAddr, consulPort,
		func() (bool, error) { return false, nil }, true)
	if err != nil {
		utils.Debug(false, "ERROR cant get ip:", err.Error())
		return
	}

	err = consul.Run()
	if err != nil {
		utils.Debug(false, "ERROR when register service ", err)
		return
	}
	defer consul.Close()

	utils.Debug(false, "✔✔✔")
	utils.Debug(false, "3. Connect to grpc servers")

	readyChan := make(chan error)
	defer close(readyChan)
	finishChan := make(chan interface{})
	defer close(finishChan)

	clients.ALL = clients.Clients{}
	clients.ALL.Init(consulAddr+consulPort, readyChan,
		finishChan, configuration.Service)

	utils.Debug(false, "✔✔✔✔")
	utils.Debug(false, "4. Launch the game lobby")

	engine.Launch(&configuration.Game, db, photo.GetImages)
	defer engine.GetLobby().Stop()

	var (
		r   = handlers.Router(db, configuration)
		srv = server.Server(r, configuration.Server, false, mainPort)
	)

	utils.Debug(false, "✔✔✔✔✔")
	utils.Debug(false, "Service", serviceName, "with id:",
		server.ServiceID(serviceName), "ready to go on",
		configuration.Server.Host+mainPort)

	server.LaunchHTTP(srv, configuration.Server, maxConn, func() {
		finishChan <- nil
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
	os.Exit(0)
}
