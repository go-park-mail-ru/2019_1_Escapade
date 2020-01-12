package server

import (
	"fmt"
	"net/http"
	"google.golang.org/grpc"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server/load_balancer"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server/service_discovery"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// ServicesI interface of service
// 	 With the Run() function the service run all it's dependencies:
//    connections to databases, another services and so on, Also
//    there it must initialize Args.Handler, that represent runnable
//    server object. If it returns an error, the server startup
//    function will stop executing the steps and the program will
//    exit with the os command.Exit().
//
//  With the Close() function, the service closes connections to other
//   services, databases, stops running gorutins, frees resources, and
//   so on. It also can return error, As well as Run(), Close() can
//   return an error, which will terminate the program with an error code
type ServiceI interface {
	Run(args *Args) error
	Router() http.Handler
	Close() error
}

type Args struct {
	Name string
	Port int

	Subnet string
	DiscoveryAddr string

	Stages []func(*Args) int
	
	Loader  ConfigutaionLoaderI
	Discovery  service_discovery.Interface
	Service ServiceI

	LoadBalancer load_balancer.Interface

	GRPC *grpc.Server

	Test bool
}

// Validate check args are correct
func (args *Args) Validate() bool {
	return re.NoNil(args, args.Loader,
		args.Discovery, args.Service) == nil
}

const (
	NOERROR       = 0
	ARGSERROR     = 1
	INPUTERROR    = 2
	CONFIGERROR   = 3
	CONSULERROR   = 4
	RUNERROR      = 5
	MAINRUNERROR  = 6
	STOPERROR     = 7
	EXTRARROR     = 8
	NOSERVERERROR = 9
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

	if !args.Validate() {
		panic(synced.Exit{Code: ARGSERROR})
	}

	stages := []func(*Args) int{}
	if args.Stages != nil {
		stages = append(stages, args.Stages...)
	}
	stages = append(stages, load, orchestration,
		runDependencies, runServer, stopDependencies)

	var errorCode = runStages(args, stages...)

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
func orchestration(args *Args) int {
	var (
			conf = args.Loader.Get().Server
			name = conf.Name
			port = args.Port
			host = GetIP(&args.Subnet)
			ttl = conf.Timeouts.TTL.Duration
			maxconn = conf.MaxConn
			addr = args.DiscoveryAddr
			check = func() (bool, error) { return false, nil }
	)
	input := new(service_discovery.Input).Init(name, port, host, 
		ttl, maxconn, addr, check)

	if (args.LoadBalancer != nil) {
		input.AddLoadBalancer(args.LoadBalancer)
	}

	utils.Debug(false, "args.Discovery.Init!")
	err := args.Discovery.Init(input).Run()
	if err != nil {
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

// run server
func runServer(args *Args) int {
	var (
		c    = args.Loader.Get().Server
		port = PortString(args.Port)
		err  error
	)

	utils.Debug(false, "Service", c.Name, "with id:",
		args.Discovery.Data().ID, "ready to go on", GetIP(&args.Subnet)+port)

	if args.GRPC != nil {
		err = LaunchGRPC(args.GRPC, c, port, func() { utils.Debug(false, "✗✗✗ Exit ✗✗✗") })
	} else {
		var srv = ConfigureServer(args.Service.Router(), c, port)
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
	args.Discovery.Close()
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

func PortString(port int) string {
	return  ":" + utils.String(port)
}