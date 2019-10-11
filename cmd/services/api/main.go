package main

import (
	"net"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	e_server "github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	_ "github.com/go-park-mail-ru/2019_1_Escapade/docs"

	"os"
)

// to generate docs, call from root "swag init -g api/main.go"

// @title Escapade Explosion API
// @version 1.0
// @description We don't have a public API, so instead of a real host(explosion.team) we specify localhost:3001. To test the following methods, git clone https://github.com/go-park-mail-ru/2019_1_Escapade, enter the root directory and run 'docker-compose up -d'

// @host localhost:3001
// @BasePath /api
func main() {
	var (
		configuration *config.Configuration
		API           *api.Handler
		err           error
	)

	utils.Debug(false, "1. Check command line arguments")

	if len(os.Args) < 6 {
		utils.Debug(false, "ERROR. Api service need 5 command line arguments. But",
			len(os.Args)-1, "get.")
		return
	}

	var (
		configurationPath = os.Args[1]
		photoPublicPath   = os.Args[2]
		photoPrivatePath  = os.Args[3]
		mainPort          = os.Args[4]
		consulPort        = os.Args[5]
		mainPortInt       int
	)

	mainPort, mainPortInt, err = e_server.Port(mainPort)
	if err != nil {
		utils.Debug(false, "ERROR - invalid server port(cant convert to int):", err.Error())
		return
	}
	consulPort = e_server.FixPort(consulPort)

	utils.Debug(false, "✔")
	utils.Debug(false, "2. Setting the environment")

	configuration, err = config.Init(configurationPath)
	if err != nil {
		utils.Debug(false, "ERROR with main configuration:", err.Error())
		return
	}

	err = photo.Init(photoPublicPath, photoPrivatePath)
	if err != nil {
		utils.Debug(false, "ERROR with photo configuration:", err.Error())
		return
	}

	API, err = api.GetHandler(configuration)
	if err != nil {
		utils.Debug(false, "ERROR with photo configuration:", err.Error())
		return
	}

	metrics.InitApi()

	//API.RandomUsers(10) // create 10 users for tests

	utils.Debug(false, "✔✔")
	utils.Debug(false, "3. Set the settings of our server and associate it with third-party")

	r := e_server.APIRouter(API, configuration.Cors, configuration.Cookie,
		configuration.Auth, configuration.AuthClient)

	srv := e_server.Server(r, configuration.Server, true, mainPort)

	// в конфиг
	var (
		serviceName = "api"
		ttl         = time.Second * 10
	)

	// /sbin/ip route|awk ' { print $7 }'

	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = configuration.Server.Host
	}

	finishHealthCheck := make(chan interface{}, 1)
	serviceID := e_server.ServiceID(serviceName)
	var serviceAddr net.IP
	serviceAddr, err = e_server.GetIP()

	newTags := []string{"api", "traefik.enable=true",
		"traefik.frontend.entryPoints=http",
		"traefik.frontend.rule=Host:api.consul.localhost"}
	consul, err := e_server.ConsulClient(serviceAddr.String(),
		serviceName, consulAddr,
		serviceID, mainPortInt, newTags, consulPort, ttl,
		func() (bool, error) { return false, nil }, finishHealthCheck)
	if err != nil {
		close(finishHealthCheck)
		utils.Debug(false, "ERROR while connecting to consul")
	}

	utils.Debug(false, "✔✔✔")
	utils.Debug(false, "Service", serviceName, "with id:", serviceID, "ready to go on", configuration.Server.Host+mainPort)

	e_server.LaunchHTTP(srv, configuration.Server, func() {
		finishHealthCheck <- nil
		API.Close()
		err := consul.Agent().ServiceDeregister(serviceID)
		if err != nil {
			utils.Debug(false, "Consul error while deregistering:", err.Error())
			return
		}
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
	os.Exit(0)
}
