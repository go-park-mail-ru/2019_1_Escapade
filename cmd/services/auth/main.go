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
	erydb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"gopkg.in/oauth2.v3/models"
)

/*
curl -H Host:api.2019-1-escapade.docker.localhost http://127.0.0.1/api/user
*/
func main() {

	var (
		configuration *config.Configuration
		err           error
		// database with users
		db        *database.DataBase
		anotherDB *erydb.DB
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
	utils.Debug(false, "3. Connect to user information storages")

	db, err = database.Init(configuration.DataBase)
	if err != nil {
		utils.Debug(false, "ERROR with user database:", err.Error())
		return
	}
	defer db.Db.Close()

	anotherDB, err = erydb.Init("postgres://eryuser:nopassword@pg-ery:5432/erybase?sslmode=disable",
		20, 20, time.Hour)
	if err != nil {
		utils.Debug(false, "ERROR with ery database:", err.Error())
		return
	}
	defer anotherDB.Close()

	utils.Debug(false, "✔✔✔")
	utils.Debug(false, "4. Connect to tokens and clients storage")

	// TODO в конфиг
	var (
		accessTokenExp    = time.Hour * 24 * 14
		refreshTokenExp   = time.Hour * 24 * 14
		isGenerateRefresh = true
		jwtSecret         = "00000000"
		link              = configuration.DataBase.AuthConnectionString

		clients = make([]*models.Client, 1)
	)
	clients[0] = &models.Client{
		ID:     "1",
		Secret: "1",
		Domain: "api.consul.localhost",
	}

	manager, tokenStore, err := e_oauth.Init(accessTokenExp, refreshTokenExp,
		isGenerateRefresh, jwtSecret, link, clients)
	if err != nil {
		utils.Debug(false, "ERROR with oauth2 equipment", err.Error())
		return
	}
	defer tokenStore.Close()

	utils.Debug(false, "✔✔✔✔")
	utils.Debug(false, "5. Set the settings of server and register in Consul")

	srv := e_oauth.Server(db, anotherDB, manager)
	r := a_handlers.Router(srv, tokenStore)
	server := e_server.Server(r, configuration.Server, true, mainPort)

	var (
		serviceName = "auth"
		ttl         = time.Second * 10
		maxConn     = 10
	)

	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = configuration.Server.Host
	}

	newTags := []string{"auth", "traefik.frontend.entryPoints=http",
		"traefik.frontend.rule=Host:auth.consul.localhost"}
	consul, err := e_server.InitConsulService(serviceName,
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

	utils.Debug(false, "✔✔✔✔✔")
	utils.Debug(false, "Service", serviceName, "with id:",
		e_server.ServiceID(serviceName), "ready to go on",
		configuration.Server.Host+mainPort)

	e_server.LaunchHTTP(server, configuration.Server, maxConn, func() {
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
}
