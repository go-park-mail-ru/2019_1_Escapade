package main

import (
	"os"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database"
	game "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/service"
)

const ARGSLEN = 7

func main() {
	args := &server.Args{
		Input:  generateInput(),
		Loader: generateLoader(),
		Consul: new(server.ConsulService),
		Service: &game.Service{
			Chat:     new(clients.Chat),
			Constant: new(constants.RepositoryFS),
			Database: new(database.Input).InitAsPSQL(),
		},
	}

	server.Run(args)
}

func generateInput() *server.Input {
	var input = new(server.Input).InitAsCMD(os.Args[6], ARGSLEN)
	input.CallInit = func() {
		input.Data.FieldPath = os.Args[4]
		input.Data.RoomPath = os.Args[5]
		input.Data.MainPort = os.Args[6]
	}
	return input
}

func generateLoader() *server.Loader {
	var loader = new(server.Loader).InitAsFS(os.Args[1])
	loader.CallExtra = func() error {
		return loader.LoadPhoto(os.Args[2], os.Args[3])
	}
	return loader
}
