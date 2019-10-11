package main

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
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

	var (
		db *database.DataBase
	)

	if db, err = database.Init(configuration.DataBase); err != nil {
		return
	}

	service := chat.NewService(db.Db, mainPortInt)

	var (
		serviceName = "chat"
		ttl         = time.Second * 10
	)

	finishHealthCheck := make(chan interface{}, 1)

	utils.Debug(false, "✔✔")
	utils.Debug(false, "3. Register in Consul")

	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = configuration.Server.Host
	}

	utils.Debug(false, "consulAddr", consulAddr)

	serviceID := server.ServiceID(serviceName)
	consul, err := server.ConsulClient(serviceName, consulAddr,
		serviceID, mainPortInt, []string{"chat"}, consulPort, ttl,
		service.Check, finishHealthCheck)

	if err != nil {
		close(finishHealthCheck)
		return
	}

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, service)

	utils.Debug(false, "✔✔✔")
	utils.Debug(false, "Service", serviceName, "with id:", serviceID, "ready to go on", configuration.Server.Host+mainPort)

	server.LaunchGRPC(grpcServer, mainPort, func() {
		finishHealthCheck <- nil
		service.Close()
		err := consul.Agent().ServiceDeregister(serviceID)
		if err != nil {
			utils.Debug(false, "Consul error while deregistering:", err.Error())
			return
		}
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
	os.Exit(0)
}
