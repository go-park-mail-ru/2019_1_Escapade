package main

import (
	"flag"
	"log"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/router"
	sd "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/service_discovery"
	lb "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/load_balancer"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/loader"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"

	factory "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/factory"

	// dont delete it for correct easyjson work
	_ "github.com/go-park-mail-ru/2019_1_Escapade/docs"
	_ "github.com/mailru/easyjson/gen"
)

// to generate docs, call from root "swag init -g api/main.go"

// @title Escapade Explosion API
// @version 1.0
// @description We don't have a public API, so instead of a real host(explosion.team) we specify localhost:3001. To test the following methods, git clone https://github.com/go-park-mail-ru/2019_1_Escapade, enter the root directory and run 'docker-compose up -d'

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://localhost:3003/auth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @host virtserver.swaggerhub.com/SmartPhoneJava/explosion/1.0.0
// @BasePath /api

var (
	port                                     int
	name                                     string
	pathToConfig, pathToSecrets, pathToPhoto string
	subnet, consulHost                       string
	CheckHTTPTimeout, checkHTTPInterval      string
)

func init() {
	flag.IntVar(&port, "port", 80, "port number(default: 80)")
	flag.StringVar(&name, "name", "api", "the service name(default: api)")
	flag.StringVar(&pathToConfig, "config", "-", "path to service configuration file")
	flag.StringVar(&pathToSecrets, "secrets", "-", "path to secrets")
	flag.StringVar(&pathToPhoto, "photo", "-", "path to configuration file of photo service")
	flag.StringVar(&subnet, "subnet", ".", "first 3 bytes of network(example: 10.10.8.)")
	flag.StringVar(&consulHost, "consul", "consul:8500", "address of consul(default: consul:8500)")
	flag.StringVar(&CheckHTTPTimeout, "http.check.timeout", "2s", "timeout of http check(default: 2s)")
	flag.StringVar(&checkHTTPInterval, "http.check.interval", "10s", "interval of http check(default: 10s)")
}

func main() {
	flag.Parse()

	var (
		conf   = loadConfigFromFS()                     // load configuration from file system
		Router = router.NewMuxRouter()                  // route paths by gorilla mux
		psql   = database.NewPostgresSQL(conf.DataBase) // use Postgresql as Database

		// handle connections by our API's handler
		Handler = factory.NewHandler(conf, psql, time.Minute, Router, subnet)
		// use Traefik as reverse proxy and Consul as Service Discovery
		balancing = traefikWithConsul(conf)

		// run and stop connection to DB and consul TTL
		// runGoroutines  = func() error {
		// 	synced.Run(context.Background(), prepareTimeout,
		// 	func() error {return psql.Open(conf.DataBase)},
		// 	balancing.Run)
		// }
		// stopGoroutines = func() error { return re.Close(psql, balancing) }

		// Create http server
		srv = server.NewHTTPServer(conf.Server, Handler, server.PortString(port))
	)

	srv.AddDependencies(psql, balancing).Run()
}

func loadConfigFromFS() *config.Configuration {
	var load = loader.NewLoader(config.NewRepositoryFS(), pathToConfig)

	if err := load.Load(); err != nil {
		log.Fatal("no main configuration found:", err.Error())
	}
	if err := load.LoadPhoto(pathToPhoto, pathToSecrets); err != nil {
		log.Fatal("no photo configuration found:", err.Error())
	}

	return load.Get()
}

func traefikWithConsul(c *config.Configuration) sd.Interface {
	var input = sd.NewInput(c.Server.Name, port,
		server.GetIP(&subnet),
		c.Server.Timeouts.TTL.Duration, c.Server.MaxConn,
		consulHost, func() (bool, error) { return false, nil })

	traefik := lb.NewTraefik(lb.NewInputEnv(name, port))

	input.AddLoadBalancer(traefik)
	var consul = sd.NewConsul(input)
	consul.SetExtraChecks(consul.HTTPCheck("http", "/api/health",
		CheckHTTPTimeout, checkHTTPInterval))
	return consul
}
