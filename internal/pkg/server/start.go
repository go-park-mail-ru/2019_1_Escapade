package server

import (
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

type CommandLineArgs struct {
	ConfigurationPath string
	PhotoPublicPath   string
	PhotoPrivatePath  string
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
	Photo bool
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
