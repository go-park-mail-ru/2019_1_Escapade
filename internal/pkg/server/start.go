package server

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	"google.golang.org/grpc"
)

type ServicesI interface {
	Run(args *Args) error
	Close() error
}

type HandlerI interface {
	Router() http.Handler
}

type Args struct {
	Input         InputI
	Loader        ConfigutaionLoaderI
	ConsulService ConsulServiceI
	Service       ServicesI
	Handler       HandlerI

	GRPC *grpc.Server

	Test bool
}

const (
	NOERROR      = 0
	ARGSERROR    = 1
	INPUTERROR   = 2
	CONFIGERROR  = 3
	CONSULERROR  = 4
	RUNERROR     = 5
	MAINRUNERROR = 6
	STOPERROR    = 7
	EXTRARROR    = 8
)

/*
Run performs all stages of loading and starting the server
 1. loading input parameters
 2. loading configuration files
 3. registration in the discovery service
 4. server startup(at this point, the execution thread is blocked because
	the server starts listening for incoming connections)
 5. Stopping the server with resource cleanup

 if an error occurs at one of the stages, a panic will be triggered. It
 will be intercepted by the synced.Exit, which will call os.Exit(code)
*/
func Run(args *Args) {
	synced.HandleExit()

	if args == nil || args.Input == nil || args.Loader == nil ||
		args.ConsulService == nil || args.Service == nil {
		panic(synced.Exit{Code: ARGSERROR})
	}

	var errorCode = runStages(args, input, load, consul, runDependencies,
		runServer, stopDependencies)

	fmt.Println("errorCode:", errorCode)
	if errorCode != NOERROR {
		panic(synced.Exit{Code: errorCode})
	}
}

// runStages every stage(func taken Args and returning code error)
func runStages(args *Args, stages ...func(*Args) int) int {
	var code = NOERROR
	for i, action := range stages {
		if code = action(args); code != NOERROR {
			printFAIL(i)
			break
		}
		printOK(i)
	}
	return code
}

// loading input parameters
func input(args *Args) int {
	if err := args.Input.CheckBefore(); err != nil {
		utils.Debug(false, "ERROR with check before init:", err.Error())
		return INPUTERROR
	}

	args.Input.Init()

	if err := args.Input.CheckAfter(); err != nil {
		utils.Debug(false, "ERROR with  check after init:", err.Error())
		return INPUTERROR
	}

	// if err := args.Input.Extra(); err != nil {
	// 	utils.Debug(false, "ERROR with input extra action:", err.Error())
	// 	return EXTRARROR
	// }

	return NOERROR
}

//  loading configuration files
func load(args *Args) int {
	if err := args.Loader.Load(); err != nil {
		utils.Debug(false, "ERROR with configuration:", err.Error())
		return CONFIGERROR
	}

	if err := args.Loader.Extra(); err != nil {
		utils.Debug(false, "ERROR with configuration extra action:", err.Error())
		return EXTRARROR
	}

	return NOERROR
}

// registration in the discovery service
func consul(args *Args) int {
	input := new(ConsulInput).Init(args.Input, args.Loader)
	consul := args.ConsulService.Init(input)

	if err := consul.Run(); err != nil {
		utils.Debug(false, "ERROR with consul:", err.Error())
		return CONSULERROR
	}
	return NOERROR
}

// run depencies
func runDependencies(args *Args) int {
	if err := args.Service.Run(args); err != nil {
		utils.Debug(false, "ERROR with running server:", err.Error())
		return RUNERROR
	}
	return NOERROR
}

// // run server
func runServer(args *Args) int {
	if args.Handler == nil {
		return MAINRUNERROR
	}

	var c = args.Loader.Get().Server
	var port = args.Input.GetData().MainPort
	var err error

	utils.Debug(false, "Service", c.Name, "with id:",
		args.ConsulService.ServiceID(), "ready to go on", GetIP()+port)

	if args.GRPC != nil {
		err = LaunchGRPC(args.GRPC, c, port, func() { utils.Debug(false, "✗✗✗ Exit ✗✗✗") })
	} else {
		var srv = ConfigureServer(args.Handler.Router(), c, port)
		err = LaunchHTTP(srv, c, func() { utils.Debug(false, "✗✗✗ Exit ✗✗✗") })
	}
	if err != nil {
		return MAINRUNERROR
	}

	return NOERROR
}

// stop depencies
func stopDependencies(args *Args) int {
	if err := args.Service.Close(); err != nil {
		utils.Debug(false, "ERROR with stopping server:", err.Error())
		return RUNERROR
	}
	args.ConsulService.Close()
	return NOERROR
}

func printOK(i int) {
	str := ""
	for a := 0; a < i+1; a++ {
		str += "✔"
	}
}

func printFAIL(i int) {
	str := ""
	for a := 0; a < i+1; a++ {
		str += "✕"
	}
}
