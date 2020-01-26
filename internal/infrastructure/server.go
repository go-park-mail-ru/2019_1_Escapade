package infrastructure

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
)

const (
	PREPAREERROR = 1
	RUNERROR     = 2
	CLOSEERROR   = 3
)

// DependencyI interface of service
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

type ServerI interface {
	Run()
	AddDependencies(dependencies ...DependencyI) ServerI
}

type ServerBase struct {
	run          func() error
	dependencies []DependencyI
	config       config.Server
}

func (server *ServerBase) Init(run func() error, config config.Server) ServerI {
	server.run = run
	server.config = config
	return server
}

func (server *ServerBase) AddDependencies(dependencies ...DependencyI) ServerI {
	server.dependencies = dependencies
	return server
}

func (server *ServerBase) Run() {
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

func (server *ServerBase) runDependencies() error {
	var actions = make([]func() error, 0)
	for _, dependency := range server.dependencies {
		actions = append(actions, dependency.Run)
	}
	timeout := server.config.Timeouts.Prepare.Duration
	return synced.Run(context.Background(), timeout, actions...)
}

func (server *ServerBase) closeDependencies() error {
	var actions = make([]func() error, 0)
	for _, dependency := range server.dependencies {
		actions = append(actions, dependency.Close)
	}
	return synced.Close(actions...)
}
