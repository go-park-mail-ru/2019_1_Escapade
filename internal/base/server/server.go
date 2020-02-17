package server

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
)

type ServerAddr struct{}

func (ServerAddr) IP(subnet *string) string {
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
}

type ServerBase struct {
	prepareDuration time.Duration
	run             func() error
	dependencies    []infrastructure.Dependency
}

func New(
	prepareDuration time.Duration,
	run func() error,
) *ServerBase {
	return &ServerBase{
		prepareDuration: prepareDuration,
		run:             run,
	}

}

func (server *ServerBase) AddDependencies(
	dependencies ...infrastructure.Dependency,
) infrastructure.Server {
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
	timeout := server.prepareDuration //server.config.Timeouts.Prepare.Duration
	return synced.Run(
		context.Background(),
		timeout,
		actions...,
	)
}

func (server *ServerBase) closeDependencies() error {
	var actions = make([]func() error, 0)
	for _, dependency := range server.dependencies {
		actions = append(actions, dependency.Close)
	}
	return synced.Close(actions...)
}
