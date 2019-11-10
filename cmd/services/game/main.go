package main

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	start "github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"os"
)

func main() {

	cla, err := start.GetCommandLineArgs(7, func() *start.CommandLineArgs {
		return &start.CommandLineArgs{
			ConfigurationPath: os.Args[1],
			PhotoPublicPath:   os.Args[2],
			PhotoPrivatePath:  os.Args[3],
			FieldPath:         os.Args[4],
			RoomPath:          os.Args[5],
			MainPort:          os.Args[6],
		}
	})
	if err != nil {
		utils.Debug(false, "ERROR with command line args", err.Error())
		return
	}

	ca := &start.ConfigurationArgs{
		HandlersMetrics: true,
		GameMetrics:     true,
		Photo:           true,
		Room:            true,
		Field:           true,
	}
	// second step
	configuration, err := start.GetConfiguration(cla, ca)
	if err != nil {
		utils.Debug(false, "ERROR with configuration", err.Error())
		return
	}

	lastArgs := &start.AllArgs{
		C:           configuration,
		CLA:         cla,
		IsWebsocket: true,
	}
	// third step
	consul := start.RegisterInConsul(lastArgs)
	// start connection to Consul
	err = consul.Run()
	if err != nil {
		utils.Debug(false, "ERROR with connection to Consul:", err.Error())
		return
	}
	defer consul.Close()

	utils.Debug(false, "3.5. Connect to grpc servers and database")

	var chatService = clients.Chat{}
	err = chatService.Init(consul, configuration.Required)
	if err != nil {
		utils.Debug(false, "ERROR with grpc connection:", err.Error())
		return
	}
	defer chatService.Close()

	utils.Debug(false, "✔✔")

	// start connection to database inside handlers
	var handler handlers.GameHandler
	err = handler.InitWithPostgresql(chatService, configuration)
	if err != nil {
		utils.Debug(false, "Database error:", err.Error())
		return
	}
	defer handler.Close()

	// forth step
	var srv = start.ConfigureServer(handler.Router(), lastArgs)

	utils.Debug(false, "Service", consul.Name, "with id:", consul.ID, "ready to go on",
		start.GetIP()+cla.MainPort)

	fmt.Println("more conf:", configuration.Game.Lobby.Intervals, configuration.Game.Lobby.ConnectionTimeout,
		configuration.Server.MaxConn)
	server.LaunchHTTP(srv, configuration.Server, func() {
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
	os.Exit(0)
}
