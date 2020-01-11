package main

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/service"

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

const ARGSLEN = 5

func main() {
	server.Run(&server.Args{
		Input:  input(),
		Loader: loader(),
		Consul: consul(),
		Service: service(),
	})
}

func input() *server.Input {
	return new(server.Input).InitAsCMD(
		server.OSArg(4), ARGSLEN)
}

func loader() *server.Loader {
	var loader = new(server.Loader).InitAsFS(server.OSArg(1))
	loader.CallExtra = func() error {
		return loader.LoadPhoto(server.OSArg(2), server.OSArg(3))
	}
	return loader
}

func consul() *server.ConsulService {
	var cs = new(server.ConsulService)
	cs.AddHTTPCheck("http","/health")
	return cs
}

func service() server.ServiceI {
	return new(api.Service).Init(
		new(database.Input).InitAsPSQL())
}

// 120 -> 62 -> 93 -> 71 -> 64
