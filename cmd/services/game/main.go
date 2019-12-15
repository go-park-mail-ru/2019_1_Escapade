package main

import (
	"os"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	pkgServer "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database"
	game "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/service"
)

const ARGSLEN = 7

func main() {

	args := &pkgServer.Args{
		Input:         generateInput(),
		Loader:        generateLoader(),
		ConsulService: new(pkgServer.ConsulService),
		Service: &game.Service{
			Chat:     new(clients.Chat),
			Consul:   new(pkgServer.ConsulService),
			Constant: new(constants.RepositoryFS),
			Database: new(database.Input).InitAsPSQL(),
		},
	}

	pkgServer.Run(args)
}

func generateInput() *pkgServer.Input {
	var input = new(pkgServer.Input)

	input.CallInit = func() {
		input.Data.FieldPath = os.Args[4]
		input.Data.RoomPath = os.Args[5]
		input.Data.MainPort = os.Args[6]
	}

	input.CallCheckBefore = func() error {
		return input.CheckBeforeDefault(ARGSLEN)
	}

	input.CallCheckAfter = func() error {
		return input.CheckAfterDefault()
	}
	return input
}

func generateLoader() *pkgServer.Loader {
	var loader = new(pkgServer.Loader)
	loader.Init(new(config.RepositoryFS), os.Args[1])
	loader.CallExtra = func() error {
		return loader.LoadPhoto(os.Args[2], os.Args[3])
	}
	return loader
}
