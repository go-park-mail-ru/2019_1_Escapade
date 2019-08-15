package main

import (
	"os"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	e_server "github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	a_handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/handlers"
	e_oauth "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/oauth"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"gopkg.in/oauth2.v3/models"
)

func main() {

	var (
		configuration *config.Configuration
		err           error
		// database with users
		db *database.DataBase
	)

	utils.Debug(false, "1. Check command line arguments")

	if len(os.Args) < 5 {
		utils.Debug(false, "ERROR. Auth service need 4 command line arguments. But",
			len(os.Args)-1, "get.")
		return
	}

	var (
		configurationPath = os.Args[1]
		mainPort          = os.Args[2]
		consulPort        = os.Args[3]
		mainPortInt       int
		// database with clients and tokens
		postgresqlPort int
	)
	postgresqlPort, err = strconv.Atoi(os.Args[4])
	if err != nil {
		utils.Debug(false, "ERROR. Wrong port of auth models database", err.Error())
		return
	}

	utils.Debug(false, "Warning: postgresqlPort declared and not used", postgresqlPort)

	mainPort, mainPortInt, err = e_server.Port(mainPort)
	if err != nil {
		utils.Debug(false, "ERROR - invalid server port(cant convert to int):", err.Error())
		return
	}
	consulPort = e_server.FixPort(consulPort)

	utils.Debug(false, "✔")
	utils.Debug(false, "2. Set the configuration")

	configuration, err = config.Init(configurationPath)
	if err != nil {
		utils.Debug(false, "ERROR with main configuration:", err.Error())
		return
	}

	utils.Debug(false, "✔✔")
	utils.Debug(false, "3. Connect to user information storage")

	db, err = database.Init(configuration.DataBase)
	if err != nil {
		utils.Debug(false, "ERROR with database:", err.Error())
		return
	}

	utils.Debug(false, "✔✔✔")
	utils.Debug(false, "4. Connect to tokens and clients storage")

	// TODO в конфиг
	var (
		accessTokenExp    = time.Hour * 2
		refreshTokenExp   = time.Hour * 24 * 14
		isGenerateRefresh = true
		jwtSecret         = "00000000"
		link              = configuration.DataBase.AuthConnectionString

		clients = make([]*models.Client, 1)
	)
	clients[0] = &models.Client{
		ID:     "1",
		Secret: "1",
		Domain: configuration.Server.Host + ":3001",
	}

	manager, tokenStore, err := e_oauth.Init(accessTokenExp, refreshTokenExp,
		isGenerateRefresh, jwtSecret, link, clients)
	if err != nil {
		utils.Debug(false, "ERROR with oauth2 equipment", err.Error())
		db.Db.Close()
		return
	}

	utils.Debug(false, "✔✔✔✔")
	utils.Debug(false, "5. Set the settings of server and register in Consul")

	srv := e_oauth.Server(db, manager)
	r := a_handlers.Router(srv, tokenStore)
	server := e_server.Server(r, configuration.Server, true, mainPort)

	var (
		serviceName = "auth"
		ttl         = time.Second * 10
	)

	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = configuration.Server.Host
	}

	finishHealthCheck := make(chan interface{}, 1)

	consul, serviceID, err := e_server.ConsulClient(serviceName, consulAddr,
		mainPort, mainPortInt, consulPort, ttl, func() (bool, error) { return false, nil },
		finishHealthCheck)
	if err != nil {
		utils.Debug(false, "ERROR while connecting to consul")
		db.Db.Close()
		tokenStore.Close()
	}

	utils.Debug(false, "✔✔✔✔✔")
	utils.Debug(false, "Service", serviceName, "with id:", serviceID, "ready to go on", configuration.Server.Host+mainPort)

	e_server.LaunchHTTP(server, configuration.Server, func() {
		finishHealthCheck <- nil
		db.Db.Close()
		tokenStore.Close()
		err := consul.Agent().ServiceDeregister(serviceID)
		if err != nil {
			utils.Debug(false, "Consul error while deregistering:", err.Error())
			return
		}
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})

}
