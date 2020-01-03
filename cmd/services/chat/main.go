package main

import (
	"os"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/database"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/service"
)

const ARGSLEN = 2

func main() {
	args := &server.Args{
		Input:  new(server.Input).InitAsCMD(os.Args[2], ARGSLEN),
		Loader: new(server.Loader).InitAsFS(os.Args[1]),
		Consul: new(server.ConsulService),
		Service: &chat.Service{
			Database: new(database.Input).InitAsPSQL(),
		},
	}
	server.Run(args)
}

// 70 -> 26
