package main



import (
	"flag"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server/load_balancer"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server/service_discovery"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/service"

	consulapi "github.com/hashicorp/consul/api"

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
	CheckHTTPTimeout, checkHTTPInterval string
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

	server.Run(&server.Args{
		Name: name,
		Port: port,

		Subnet:        subnet,
		DiscoveryAddr: consulHost,

		Loader: loader(),
		Discovery: &service_discovery.Consul{
			ExtraChecks: func(c *service_discovery.Consul) []*consulapi.AgentServiceCheck {
				return []*consulapi.AgentServiceCheck{
					c.HTTPCheck("http", "/api/health", 
						CheckHTTPTimeout, checkHTTPInterval),
				}
			},
		},

		Service: new(api.Service).Init(subnet,
			new(database.Input).InitAsPSQL()),

		LoadBalancer: new(load_balancer.Traefik).Init(
			new(load_balancer.InputEnv).Init(name, port)),
	})
}

func loader() *server.Loader {
	var loader = new(server.Loader).InitAsFS(pathToConfig)
	loader.CallExtra = func() error {
		return loader.LoadPhoto(pathToPhoto, pathToSecrets)
	}
	return loader
}

// 120 -> 62 -> 93 -> 71 -> 64 -> 81 -> 92
