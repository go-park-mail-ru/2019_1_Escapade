package main
/*
import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/metrics"
	erydatabase "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"
	eryhandlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/handlers"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
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
		H             *eryhandlers.Handler
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

	mainPort, mainPortInt, err = server.Port(mainPort)
	if err != nil {
		utils.Debug(false, "ERROR - invalid server port(cant convert to int):", err.Error())
		return
	}
	consulPort = server.FixPort(consulPort)

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

	db, err := erydatabase.Init("postgres://eryuser:nopassword@pg-ery:5432/erybase?sslmode=disable",
		20, 20, time.Hour)
	if err != nil {
		utils.Debug(false, "ERROR with database:", err.Error())
		return
	}

	H = eryhandlers.Init(db, configuration)
	if err != nil {
		utils.Debug(false, "ERROR with photo configuration:", err.Error())
		return
	}

	metrics.Init()

	//API.RandomUsers(10) // create 10 users for tests

	utils.Debug(false, "✔✔")
	utils.Debug(false, "3. Set the settings of our server and associate it with third-party")

	r := eryhandlers.Router(H, configuration.Cors, configuration.Cookie,
		configuration.Auth, configuration.AuthClient)

	srv := server.Server(r, configuration.Server, true, mainPort)

	// в конфиг
	var (
		serviceName = "ery"
		ttl         = time.Second * 10
		maxConn     = 100
	)

	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = configuration.Server.Host
	}

	finishHealthCheck := make(chan interface{}, 1)
	serviceID := server.ServiceID(serviceName)
	consul, err := server.ConsulClient("127.0.0.1", serviceName, consulAddr,
		serviceID, mainPortInt, []string{"ery"}, consulPort, ttl,
		func() (bool, error) { return false, nil }, finishHealthCheck)

	if err != nil {
		close(finishHealthCheck)
		utils.Debug(false, "ERROR while connecting to consul")
	}

	utils.Debug(false, "✔✔✔")
	utils.Debug(false, "Service", serviceName, "with id:", serviceID, "ready to go on", configuration.Server.Host+mainPort)

	server.LaunchHTTP(srv, configuration.Server, maxConn, func() {
		finishHealthCheck <- nil
		H.Close()
		err := consul.Agent().ServiceDeregister(serviceID)
		if err != nil {
			utils.Debug(false, "Consul error while deregistering:", err.Error())
			return
		}
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
	os.Exit(0)
}
*/