package main

import (
	"flag"
	"time"

	// "time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/auth/oauth2"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/database/postgresql"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/error/tracerr"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/loadbalancer/traefik"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/logger/log"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/metrics/prometheus"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/photo/aws"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/router/gorillamux"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/server/http"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/servicediscovery/consul"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/factory"

	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration/loader/cleanenv"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration/models"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/auth/oauth2"
	// postgresql "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/database/postresql"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/error/tracerr"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/loadbalancer/traefik"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/logger/log"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/metrics/prometheus"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/photo/aws"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/router/gorillamux"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/server/http"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/servicediscovery/consul"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"

	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/factory"

	// miauth "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/auth"
	// micors "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/cors"
	// milogger "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/logger"
	// mimetrics "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/metrics"
	// mirecover "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/recover"

	miauth "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/auth"
	micors "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/cors"
	milogger "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/logger"
	mimetrics "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/metrics"
	mirecover "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/middleware/recover"

	// dont delete it for correct easyjson work
	_ "github.com/go-park-mail-ru/2019_1_Escapade/docs"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration/loader/cleanenv"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	_ "github.com/mailru/easyjson/gen"
)

const ERROR = 1

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

func main() {
	synced.HandleExit()
	var (
		confInfrastructure, confAuth, confCors   string
		confDatabase, confPhoto, confSecretPhoto string
	)
	flag.StringVar(&confInfrastructure, "infrastructure-config", "-",
		"path to service configuration file")
	flag.StringVar(&confAuth, "auth-config", "-",
		"path to auth configuration file")
	flag.StringVar(&confCors, "cors-config", "-",
		"path to cors configuration file")
	flag.StringVar(&confDatabase, "database-config", "-",
		"path to database configuration file")
	flag.StringVar(&confPhoto, "photo-config", "-",
		"path to public photo configuration file")
	flag.StringVar(&confSecretPhoto, "photo-secret-config", "-",
		"path to secret photo configuration file")
	flag.Parse()

	loader := new(cleanenv.Loader)

	// initialize logger via log
	var logger = log.New()

	var c = models.All{}

	// load inrastructure configuration
	err := loader.Load(confInfrastructure, &c)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// load auth configuration
	err = loader.Load(confAuth, &c.Auth)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// load cors configuration
	err = loader.Load(confCors, &c.Cors)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// load database configuration
	err = loader.Load(confDatabase, &c.Database.ConnectionString)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// load photo configuration
	err = loader.Load(confPhoto, &c.Photo)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// load photo secrret configuration
	err = loader.Load(confSecretPhoto, &c.Photo)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	logger.Println("Auth", c.Auth)
	logger.Println("Cors", c.Cors)
	logger.Println("Database", c.Database)
	logger.Println("LoadBalancer", c.LoadBalancer)
	logger.Println("Photo", c.Photo)
	logger.Println("Server", c.Server)
	logger.Println("ServiceDiscovery", c.ServiceDiscovery)

	// initialize error tracer via ztrue/tracerr
	var errTrace = tracerr.New()

	// initialize metrics via prometheus
	var metrics = prometheus.New()

	// initialize auth service via oauth2
	// os.Getenv("AUTH_ADDRESS")
	auth, err := oauth2.New(&c.Auth, errTrace, logger)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// auth middleware
	mwrAuth, err := miauth.New(auth, logger)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// logger middleware
	mwrLogger := milogger.New(logger)

	// recover middleware
	mvrRecover := mirecover.New(logger)

	// metrics middleware
	mvrMetrics, err := mimetrics.New(metrics, c.ServiceDiscovery.Subnet)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// cors middleware
	mvrCors, err := micors.New(c.Cors, logger, errTrace)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	var (
		// middleware for requests that require authorization
		mwrWithAuth = []infrastructure.Middleware{
			mvrRecover, mvrCors, mwrAuth,
		}

		// middleware for requests that do not require authorization
		mwrNoneAuth = []infrastructure.Middleware{
			mvrRecover, mvrCors, mwrLogger, mvrMetrics,
		}
	)

	// initialize photo service via aws-sdk-go/aws
	photo, err := aws.New(c.Photo, logger)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// initialize database via postgresql
	db, err := postgresql.New(&c.Database)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// initialize router via gorilla/mux
	router := gorillamux.New(logger, errTrace)

	// initialize load balancer via traefik
	loadBalancer, err := traefik.New(&c.LoadBalancer)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}

	// initialize service discovery via Consul
	serviceDiscovery, err := consul.New(
		&c.ServiceDiscovery,
		func() (bool, error) { return false, nil },
		loadBalancer,
		logger,
		errTrace,
	)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}
	// add http check for our service
	serviceDiscovery.AddCheckHTTP(
		"http",
		"/api/health",
		c.ServiceDiscovery.HTTPTimeout.String(),
		c.ServiceDiscovery.HTTPInterval.String(),
	)

	handler := factory.NewHandler(
		auth,
		photo,
		logger,
		errTrace,
		db,
		router,
		mwrNoneAuth,
		mwrWithAuth,
		time.Minute, // TODO в конфиг контекстовый таймаут
	)

	server, err := http.New(&c.Server, handler, logger)
	if err != nil {
		logger.Println(err)
		panic(synced.Exit{Code: ERROR})
	}
	server.AddDependencies(db, serviceDiscovery).Run()
}

// 136 -> 181 -> 195 -> 201 -> 270
