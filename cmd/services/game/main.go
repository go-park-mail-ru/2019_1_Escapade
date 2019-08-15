package main

import (
	"time"

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

	metrics.InitApi()
	metrics.InitGame()

	utils.Debug(false, "✔✔")
	utils.Debug(false, "3. Register in consul")

	var (
		serviceName = "game"
		ttl         = time.Second * 10
	)

	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = configuration.Server.Host
	}

	finishHealthCheck := make(chan interface{}, 1)

	consul, serviceID, err := server.ConsulClient(serviceName, consulAddr,
		mainPort, mainPortInt, consulPort, ttl, func() (bool, error) { return false, nil },
		finishHealthCheck)
	if err != nil {
		utils.Debug(false, "ERROR while connecting to consul")
		db.Db.Close()
	}

	utils.Debug(false, "✔✔✔")
	utils.Debug(false, "3. Connect to grpc servers")

	readyChan := make(chan error)
	finishChan := make(chan interface{})

	clients.ALL.Init(consulAddr+consulPort, readyChan, finishChan, configuration.Services...)

	utils.Debug(false, "✔✔✔")
	utils.Debug(false, "4. Launch the game lobby")
	handler = api.Init(db, configuration)

	game.Launch(&configuration.Game, db, photo.GetImages)

	var (
		r   = server.GameRouter(handler, configuration.Cors)
		srv = server.Server(r, configuration.Server, false, mainPort)
	)

	utils.Debug(false, "✔✔✔✔")
	utils.Debug(false, "Service", serviceName, "with id:", serviceID, "ready to go on", configuration.Server.Host+mainPort)

	server.LaunchHTTP(srv, configuration.Server, func() {
		finishChan <- nil
		finishHealthCheck <- nil
		game.GetLobby().Stop()
		close(readyChan)
		close(finishChan)
		err := consul.Agent().ServiceDeregister(serviceID)
		if err != nil {
			utils.Debug(false, "Consul error while deregistering:", err.Error())
			return
		}
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
	os.Exit(0)
}
