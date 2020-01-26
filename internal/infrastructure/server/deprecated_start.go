package server
/*
import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server/serve"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	"google.golang.org/grpc"
)

const (
	PREPAREERROR = 1
	RUNERROR     = 2
	CLOSEERROR   = 3
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
type DependencyI interface {
	Run() error
	Close() error
}

type Server struct {
	run          func() error
	dependencies []DependencyI
	config       config.Server
}

func NewHTTPServer(c config.Server, handler http.Handler, port string) *Server {
	run := serve.StageRunServerHTTP(c, handler, port)
	return newServer(run, c)
}

func NewGRPCServer(c config.Server, GRPC *grpc.Server, port string) *Server {
	run := serve.StageRunServerGRPC(c, GRPC, port)
	return newServer(run, c)
}

func newServer(run func() error, c config.Server) *Server {
	return &Server{
		config: c,
		run:    run,
	}
}

func (server *Server) AddDependencies(dependencies ...DependencyI) *Server {
	server.dependencies = dependencies
	return server
}

func (server *Server) Run() {
	synced.HandleExit()

	// open connections, run extra specified goroutines
	err := server.runDependencies()
	if err != nil {
		panic(synced.Exit{Code: PREPAREERROR})
	}

	// close connections, stop specified goroutines
	defer func() {
		err = server.closeDependencies()
		if err != nil {
			panic(synced.Exit{Code: CLOSEERROR})
		}
	}()

	// run the server
	err = server.run()
	if err != nil {
		panic(synced.Exit{Code: RUNERROR})
	}
}

func (server *Server) runDependencies() error {
	var actions = make([]func() error, 0)
	for _, dependency := range server.dependencies {
		actions = append(actions, dependency.Run)
	}
	timeout := server.config.Timeouts.Prepare.Duration
	return synced.Run(context.Background(), timeout, actions...)
}

func (server *Server) closeDependencies() error {
	var actions = make([]func() error, 0)
	for _, dependency := range server.dependencies {
		actions = append(actions, dependency.Close)
	}
	return synced.Close(actions...)
}

func PortString(port int) string {
	return ":" + utils.String(port)
}

func Port(port string) (string, int, error) {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return port, 0, err
	}
	return ":" + port, intPort, err
}

func GetIP(subnet *string) string {
	var ips string
	ifaces, err := net.Interfaces()
	if err != nil {
		return err.Error()
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return err.Error()
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ipsting := ip.String()
			fmt.Println("ips:", ipsting)
			if subnet == nil {
				ips += " " + ipsting
			} else if strings.HasPrefix(ipsting, *subnet) {
				return ipsting
			}
		}
	}
	if subnet == nil {
		return ips
	}
	return "error: no networks. Change subnet!"
}*/
