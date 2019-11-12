package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	ametrics "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/metrics"
	gmetrics "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

type CommandLineArgs struct {
	ConfigurationPath string
	PhotoPublicPath   string
	PhotoPrivatePath  string
	FieldPath         string
	RoomPath          string
	MainPort          string
	MainPortInt       int
}

func GetCommandLineArgs(argsNeed int, init func() *CommandLineArgs) (*CommandLineArgs, error) {
	utils.Debug(false, "1. Check command line arguments")

	argsGet := len(os.Args)
	if argsGet < argsNeed {
		utils.Debug(false, "ERROR. Servce need ", argsNeed, " command line arguments. But",
			argsGet-1, "get.")
		return nil, re.ErrorServer()
	}

	var (
		input = init()
		err   error
	)
	input.MainPort, input.MainPortInt, err = Port(input.MainPort)
	if err != nil {
		utils.Debug(false, "ERROR - invalid server port(cant convert to int):", err.Error())
		return nil, err
	}
	utils.Debug(false, "✔")
	return input, err
}

type ConfigurationArgs struct {
	HandlersMetrics bool
	GameMetrics     bool
	Photo           bool
	Field           bool
	Room            bool
}

func GetConfiguration(cla *CommandLineArgs, ca *ConfigurationArgs) (*config.Configuration, error) {
	utils.Debug(false, "2. Setting the configuration")

	configuration, err := config.Init(cla.ConfigurationPath)
	if err != nil {
		utils.Debug(false, "ERROR with main configuration:", err.Error())
		return nil, err
	}
	if ca.Photo {
		err = photo.Init(cla.PhotoPublicPath, cla.PhotoPrivatePath)
		if err != nil {
			utils.Debug(false, "ERROR with photo configuration:", err.Error())
			return nil, err
		}
	}
	if ca.Field {
		err = constants.InitField(cla.FieldPath)
		if err != nil {
			utils.Debug(false, "Initialization error with field constants:", err.Error())
			return nil, err
		}
	}
	if ca.Room {
		err = constants.InitRoom(cla.RoomPath)
		if err != nil {
			utils.Debug(false, "Initialization error with room constants:", err.Error())
			return nil, err
		}
	}
	if ca.HandlersMetrics {
		ametrics.Init()
	}
	if ca.GameMetrics {
		gmetrics.Init()
	}
	utils.Debug(false, "✔✔")
	return configuration, nil
}

type AllArgs struct {
	CLA                *CommandLineArgs
	C                  *config.Configuration
	IsHTTPS            bool
	DisableTraefik     bool
	WithoutExecTimeout bool
}

func RegisterInConsul(aa *AllArgs) *ConsulService {
	utils.Debug(false, "3. Register the service in Consul discovery")

	var (
		consulAddr = os.Getenv("CONSUL_ADDRESS")
		name       = aa.C.Server.Name
		tags       = []string{name, "traefik.frontend.entryPoints=http",
			"traefik.frontend.rule=Host:" + name + ".consul.localhost"}
	)
	entrypoint := "http"
	if aa.IsHTTPS {
		entrypoint = "https"
	}
	fmt.Println("entrypoint:", entrypoint)
	tags = append(tags, "traefik.frontend.entryPoints="+entrypoint)

	consulInput := &ConsulInput{
		Name:          name,
		Port:          aa.CLA.MainPortInt,
		Tags:          tags,
		TTL:           aa.C.Server.Timeouts.TTL.Duration,
		MaxConn:       aa.C.Server.MaxConn,
		ConsulHost:    consulAddr,
		ConsulPort:    ":8500",
		Check:         func() (bool, error) { return false, nil },
		EnableTraefik: !aa.DisableTraefik,
	}
	utils.Debug(false, "✔✔✔")
	return InitConsulService(consulInput)
}

func ConfigureServer(handler http.Handler, aa *AllArgs) *http.Server {
	utils.Debug(false, "4. Configure server")

	serverConfig := aa.C.Server
	fmt.Println("conf:", serverConfig.Timeouts.Read.Duration, serverConfig.Timeouts.Write.Duration,
		serverConfig.Timeouts.Idle.Duration)
	var execT = serverConfig.Timeouts.Exec.Duration

	if execT > time.Duration(time.Second) && !aa.IsHTTPS && !aa.WithoutExecTimeout {
		handler = http.TimeoutHandler(handler, execT, "ESCAPADE DEBUG Timeout!")
	}

	srv := &http.Server{
		Addr:         aa.CLA.MainPort,
		ReadTimeout:  serverConfig.Timeouts.Read.Duration,
		WriteTimeout: serverConfig.Timeouts.Write.Duration,
		IdleTimeout:  serverConfig.Timeouts.Idle.Duration,
		Handler:      handler,
		// ConnState: func(n net.Conn, c http.ConnState) {
		// 	fmt.Println("--------------new conn state:", c)
		// },
		MaxHeaderBytes: serverConfig.MaxHeaderBytes,
	}
	utils.Debug(false, "✔✔✔✔")
	return srv
}
