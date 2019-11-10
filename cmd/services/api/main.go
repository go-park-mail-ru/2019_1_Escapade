package main

import (
	"os"

	_ "github.com/go-park-mail-ru/2019_1_Escapade/docs"
	start "github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	// dont delete it for correct easyjson work
	_ "github.com/mailru/easyjson/gen"
)

// to generate docs, call from root "swag init -g api/main.go"

// @title Escapade Explosion API
// @version 1.0
// @description We don't have a public API, so instead of a real host(explosion.team) we specify localhost:3001. To test the following methods, git clone https://github.com/go-park-mail-ru/2019_1_Escapade, enter the root directory and run 'docker-compose up -d'

// @host localhost:3001
// @BasePath /api
func main() {
	// first step
	cla, err := start.GetCommandLineArgs(5, func() *start.CommandLineArgs {
		return &start.CommandLineArgs{
			ConfigurationPath: os.Args[1],
			PhotoPublicPath:   os.Args[2],
			PhotoPrivatePath:  os.Args[3],
			MainPort:          os.Args[4],
		}
	})
	if err != nil {
		utils.Debug(false, "ERROR with command line args", err.Error())
		return
	}
	ca := &start.ConfigurationArgs{
		HandlersMetrics: true,
		Photo:           true,
	}
	// second step
	configuration, err := start.GetConfiguration(cla, ca)
	if err != nil {
		utils.Debug(false, "ERROR with configuration", err.Error())
		return
	}

	// start connection to database inside handlers
	var API = &api.Handlers{}
	err = API.InitWithPostgreSQL(configuration)
	if err != nil {
		utils.Debug(false, "ERROR with connection to database:", err.Error())
		return
	}
	defer API.Close()

	// third step
	consul := start.RegisterInConsul(cla, configuration)
	consul.AddHTTPCheck("http", "/health")

	// start connection to Consul
	err = consul.Run()
	if err != nil {
		utils.Debug(false, "ERROR with connection to Consul:", err.Error())
		return
	}
	defer consul.Close()

	// forth step
	server := start.ConfigureServer(API.Router(),
		configuration.Server, cla)

	utils.Debug(false, "Service", consul.Name, "with id:", consul.ID, "ready to go on",
		start.GetIP()+cla.MainPort)

	// go!
	start.LaunchHTTP(server, configuration.Server, func() {
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
}

// 120 -> 62
