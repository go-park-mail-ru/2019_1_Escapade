package main

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"google.golang.org/grpc"

	"os"
)

func main() {

	var (
		configuration *config.Configuration
		err           error
	)

	utils.Debug(false, "1. Check command line arguments")

	if len(os.Args) < 4 {
		utils.Debug(false, "ERROR. Chat service need 3 command line argument. But",
			len(os.Args)-1, "get.")
		return
	}

	var (
		configurationPath = os.Args[1]
		mainPort          = os.Args[2]
		consulPort        = os.Args[3]
		mainPortInt       int
	)

	mainPort, mainPortInt, err = server.Port(mainPort)
	if err != nil {
		utils.Debug(false, "ERROR - invalid server port(cant convert to int):", err.Error())
		return
	}
	consulPort = server.FixPort(consulPort)

	utils.Debug(false, "✔")
	utils.Debug(false, "2. Set the configuration and connect to database")

	configuration, err = config.Init(configurationPath)
	if err != nil {
		utils.Debug(false, "ERROR with main configuration:", err.Error())
		return
	}

	var handler chat.Handler
	handler.InitWithPostgreSQL(configuration)
	defer handler.Close()

	var (
		serviceName = "chat"
		ttl         = time.Second * 10
		maxConn     = 40
	)

	utils.Debug(false, "✔✔")
	utils.Debug(false, "3. Register in Consul")

	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = configuration.Server.Host
	}

	utils.Debug(false, "consulAddr", consulAddr)

	newTags := []string{"api", "traefik.frontend.entryPoints=http",
		"traefik.frontend.rule=Host:api.consul.localhost"}

	consul, err := server.InitConsulService(serviceName,
		mainPortInt, newTags, ttl, maxConn, consulAddr, consulPort,
		func() (bool, error) { return false, nil }, true)
	if err != nil {
		utils.Debug(false, "ERROR cant get ip:", err.Error())
		return
	}

	err = consul.Run()
	if err != nil {
		utils.Debug(false, "ERROR when register service ", err)
		return
	}
	defer consul.Close()

	if err != nil {
		return
	}

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &handler)

	utils.Debug(false, "✔✔✔")
	utils.Debug(false, "Service", serviceName, "with id:",
		server.ServiceID(serviceName), "ready to go on",
		configuration.Server.Host+mainPort)

	server.LaunchGRPC(grpcServer, mainPort, func() {
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
	os.Exit(0)
}
