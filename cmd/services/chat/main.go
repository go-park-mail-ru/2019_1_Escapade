package main

import (
	"google.golang.org/grpc"
	"os"

	start "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

func main() {
	synced.HandleExit()
	// first step
	cla, err := start.GetCommandLineArgs(3, func() *start.CommandLineArgs {
		return &start.CommandLineArgs{
			ConfigurationPath: os.Args[1],
			MainPort:          os.Args[2],
		}
	})
	if err != nil {
		utils.Debug(false, "ERROR with command line args", err.Error())
		panic(synced.Exit{Code: 1})
	}

	ca := &start.ConfigurationArgs{}
	// second step
	configuration, err := start.GetConfiguration(cla, ca)
	if err != nil {
		utils.Debug(false, "ERROR with configuration", err.Error())
		panic(synced.Exit{Code: 2})
	}

	var handler handlers.Handler
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
		panic(synced.Exit{Code: 3})
	}
	defer consul.Close()

	grpcServer := grpc.NewServer()
	proto.RegisterChatServiceServer(grpcServer, &handler)

	utils.Debug(false, "Service", consul.Name, "with id:", consul.ID, "ready to go on",
		start.GetIP()+cla.MainPort)

	start.LaunchGRPC(grpcServer, configuration.Server, cla.MainPort, func() {
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
}
