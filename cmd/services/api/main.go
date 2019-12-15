package main

import (
	"os"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/service"

	// dont delete it for correct easyjson work
	_ "github.com/go-park-mail-ru/2019_1_Escapade/docs/api"
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
		Input:         generateInput(),
		Loader:        generateLoader(),
		ConsulService: new(server.ConsulService),
		Service: &api.Service{
			Database: new(database.Input).InitAsPSQL(),
		},
	})
}

func generateInput() *server.Input {
	var input = new(server.Input)

	input.CallInit = func() {
		input.Data.MainPort = os.Args[4]
	}

	input.CallCheckBefore = func() error {
		return input.CheckBeforeDefault(ARGSLEN)
	}

	input.CallCheckAfter = func() error {
		return input.CheckAfterDefault()
	}
	return input
}

func generateLoader() *server.Loader {
	var loader = new(server.Loader)
	loader.Init(new(config.RepositoryFS), os.Args[1])
	loader.CallExtra = func() error {
		return loader.LoadPhoto(os.Args[2], os.Args[3])
	}
	return loader
}

// 120 -> 62 -> 93 -> 71
