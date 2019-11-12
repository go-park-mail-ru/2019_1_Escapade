package main

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	start "github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"google.golang.org/grpc"

	"os"
)

func main() {

	// first step
	cla, err := start.GetCommandLineArgs(3, func() *start.CommandLineArgs {
		return &start.CommandLineArgs{
			ConfigurationPath: os.Args[1],
			MainPort:          os.Args[2],
		}
	})
	if err != nil {
		utils.Debug(false, "ERROR with command line args", err.Error())
		return
	}

	ca := &start.ConfigurationArgs{}
	// second step
	configuration, err := start.GetConfiguration(cla, ca)
	if err != nil {
		utils.Debug(false, "ERROR with configuration", err.Error())
		return
	}

	var handler chat.Handler
	handler.InitWithPostgreSQL(configuration)
	defer handler.Close()

	lastArgs := &start.AllArgs{
		C:   configuration,
		CLA: cla,
	}
	consul := start.RegisterInConsul(lastArgs)
	err = consul.Run()
	if err != nil {
		utils.Debug(false, "ERROR with connection to Consul:", err.Error())
		return
	}
	defer consul.Close()

	grpcServer := grpc.NewServer()
	chat.RegisterChatServiceServer(grpcServer, &handler)

	utils.Debug(false, "Service", consul.Name, "with id:", consul.ID, "ready to go on",
		start.GetIP()+cla.MainPort)

	server.LaunchGRPC(grpcServer, configuration.Server, cla.MainPort, func() {
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
}
