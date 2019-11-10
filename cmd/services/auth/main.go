package main

import (
	"os"
	"time"

	start "github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	user_db "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	a_handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/handlers"
	e_oauth "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/oauth"
	ery_db "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"gopkg.in/oauth2.v3/models"
)

/*
curl -H Host:api.2019-1-escapade.docker.localhost http://127.0.0.1/api/user
*/
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

	// start connection to main database
	userDB := &user_db.UserUseCase{}
	userDB.InitDBWithSQLPQ(configuration.DataBase)
	defer userDB.Close()

	// start connection to erythocyte database
	eryDB, err := ery_db.Init("postgres://eryuser:nopassword@pg-ery:5432/erybase?sslmode=disable",
		20, 20, time.Hour)
	if err != nil {
		utils.Debug(false, "ERROR with ery database:", err.Error())
		return
	}
	defer eryDB.Close()

	clients := []*models.Client{&models.Client{
		ID:     "1",
		Secret: "1",
		Domain: "api.consul.localhost",
	}}

	// start connection to auth database
	manager, tokenStore, err := e_oauth.Init(configuration, clients)
	if err != nil {
		utils.Debug(false, "ERROR with oauth2 equipment", err.Error())
		return
	}
	defer tokenStore.Close()

	// third step
	consul := start.RegisterInConsul(cla, configuration)

	// start connection to Consul
	err = consul.Run()
	if err != nil {
		utils.Debug(false, "ERROR with connection to Consul:", err.Error())
		return
	}
	defer consul.Close()

	/// forth step
	srv := e_oauth.Server(userDB, eryDB, manager)
	r := a_handlers.Router(srv, tokenStore)
	server := start.ConfigureServer(r, configuration.Server, cla)

	utils.Debug(false, "Service", consul.Name, "with id:", consul.ID, "ready to go on",
		start.GetIP(), consul.Port)

	start.LaunchHTTP(server, configuration.Server, func() {
		utils.Debug(false, "✗✗✗ Exit ✗✗✗")
	})
}
