package grpc

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"golang.org/x/net/netutil"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/implementation/server"
)

type ServerGRPC struct {
	server.ServerBase
}

// New instance of ServerGRPC
func New(
	conf configuration.ServerRepository,
	GRPC *grpc.Server,
	logger infrastructure.Logger,
) (*ServerGRPC, error) {
	// check configuration repository given
	if conf == nil {
		return nil, errors.New(ErrNoConfiguration)
	}
	var c = conf.Get()

	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}

	// check grpc server given
	if GRPC == nil {
		return nil, errors.New(ErrNoGRPC)
	}

	return &ServerGRPC{
		ServerBase: *server.New(
			c.Timeouts.Prepare,
			func() error {
				return serveGRPC(
					GRPC,
					c,
					logger,
					func() {
						log.Println("✗✗✗ Exit ✗✗✗")
					},
				)
			},
		),
	}, nil

}

func serveGRPC(
	grpcServer *grpc.Server,
	c configuration.Server,
	log infrastructure.Logger,
	lastFunc func(),
) error {
	var (
		errChan  = make(chan error)
		stopChan = make(chan os.Signal)
	)

	defer func() {
		close(stopChan)
		close(errChan)
		lastFunc()
	}()

	l, err := net.Listen(Protocol, c.Port)

	if err != nil {
		log.Println("Listen error", err.Error())
		return err
	}

	defer l.Close()

	l = netutil.LimitListener(l, c.MaxConn)

	signal.Notify(stopChan, os.Interrupt)

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			errChan <- err
			close(errChan)
			log.Println("Serving error:", err.Error())
		}
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		log.Println("Fatal error: ", err.Error())
		return err
	case <-stopChan:
		grpcServer.GracefulStop()
	}
	return nil
}
