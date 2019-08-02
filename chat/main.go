package main

import (
	chat "github.com/go-park-mail-ru/2019_1_Escapade/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"google.golang.org/grpc"

	"net"
	"os"
)

func main() {

	var (
		configuration *config.Configuration
		err           error
	)

	if len(os.Args) < 2 {
		utils.Debug(false, "Api service need 1 command line argument. But",
			len(os.Args)-1, "get.")
		return
	}

	configurationPath := os.Args[1]

	configuration, err = config.Init(configurationPath)
	if err != nil {
		utils.Debug(false, "Initialization error with main configuration:", err.Error())
		return
	}

	var (
		db *database.DataBase
	)

	if db, err = database.Init(configuration.DataBase); err != nil {
		return
	}

	service := &chat.Service{
		DB: db.Db,
	}

	port := server.Port(configuration)

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		utils.Debug(true, "failed to listen:", err.Error())
	}
	s := grpc.NewServer()

	chat.RegisterChatServiceServer(s, service)

	server.LaunchGRPC(s, lis, func() {
		service.DB.Close()
	})
	os.Exit(0)
}
