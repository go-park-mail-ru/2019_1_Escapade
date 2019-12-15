package main

import (
	"os"

	_ "github.com/go-park-mail-ru/2019_1_Escapade/docs/auth"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/database"
	auth "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/service"
)

// to generate docs, call from root "swag init -g auth/main.go"

// @title Escapade Explosion AUTH
// @version 1.0
// @description We don't have a public API, so instead of a real host(explosion.team) we specify localhost:3003. To test the following methods, git clone https://github.com/go-park-mail-ru/2019_1_Escapade, enter the root directory and run 'docker-compose up -d'

// @host localhost:3003
// @BasePath /auth

/*
curl -H Host:api.2019-1-escapade.docker.localhost http://127.0.0.1/api/user
*/

const ARGSLEN = 3

func main() {
	server.Run(&server.Args{
		Input:         generateInput(),
		Loader:        generateLoader(),
		ConsulService: new(server.ConsulService),
		Service: &auth.Service{
			Database:    new(database.Input).InitAsPSQL(),
			RepositoryI: &clients.RepositoryHC{},
		},
	})
}

func generateInput() *server.Input {
	var input = new(server.Input)

	input.CallInit = func() {
		input.Data.MainPort = os.Args[2]
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
	return loader
}

// 111 -> 66
