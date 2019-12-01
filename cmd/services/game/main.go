package main

import (
	"fmt"
	"os"

	start "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
)

func main() {
	synced.HandleExit()
	cla, err := start.GetCommandLineArgs(7, func() *start.CommandLineArgs {
		return &start.CommandLineArgs{
			ConfigurationPath: os.Args[1],
			PhotoPublicPath:   os.Args[2],
			PhotoPrivatePath:  os.Args[3],
			MainPort:          os.Args[6],
		}
	})
	fieldPath := os.Args[4]
	roomPath := os.Args[5]
	if err != nil {
		utils.Debug(false, "ERROR with command line args", err.Error())
		panic(synced.Exit{Code: 1})
	}

	ca := &start.ConfigurationArgs{
		Photo: true,
	}
	// second step
	configuration, err := start.GetConfiguration(cla, ca)
	if err != nil {
		utils.Debug(false, "ERROR with configuration", err.Error())
		panic(synced.Exit{Code: 2})
	}

	lastArgs := &start.AllArgs{
		C:                  configuration,
		CLA:                cla,
		WithoutExecTimeout: true,
	}
	// third step
	consul := start.RegisterInConsul(lastArgs)
	// start connection to Consul
	err = consul.Run()
	if err != nil {
		utils.Debug(false, "ERROR with connection to Consul:", err.Error())
		panic(synced.Exit{Code: 3})
	}
	defer consul.Close()

	utils.Debug(false, "3.5. Connect to grpc servers and database")

	var chatService = clients.Chat{}
	err = chatService.Init(consul, configuration.Required)
	if err != nil {
		utils.Debug(false, "ERROR with grpc connection:", err.Error())
		panic(synced.Exit{Code: 4})
	}
	defer chatService.Close()

	utils.Debug(false, "✔✔")

	var gca = &handlers.ConfigurationArgs{
		C:         configuration,
		FieldPath: fieldPath,
		RoomPath:  roomPath,
	}
	// start connection to database inside handlers
	var handler handlers.GameHandler
	err = handler.InitWithPostgresql(chatService, gca)
	if err != nil {
		utils.Debug(false, "Database error:", err.Error())
		panic(synced.Exit{Code: 5})
	}
	defer handler.Close()

	// forth step
	var srv = start.ConfigureServer(handler.Router(), lastArgs)

	utils.Debug(false, "Service", consul.Name, "with id:", consul.ID, "ready to go on",
		start.GetIP()+cla.MainPort)

	fmt.Println("more conf:", configuration.Game.Lobby.Intervals, configuration.Game.Lobby.ConnectionTimeout,
		configuration.Server.MaxConn)
	start.LaunchHTTP(srv, configuration.Server, func() {
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
}
