package grpc

import (
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"golang.org/x/net/netutil"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/server"
)

type ServerGRPC struct {
	server.ServerBase

	log   infrastructure.LoggerI
	trace infrastructure.ErrorTrace
}

func New(
	c server.Configuration,
	GRPC *grpc.Server,
	log infrastructure.LoggerI,
	trace infrastructure.ErrorTrace,
) *ServerGRPC {
	if log == nil {
		log = &infrastructure.LoggerEmpty{}
	}
	if trace == nil {
		trace = &infrastructure.ErrorTraceDefault{}
	}
	return &ServerGRPC{
		ServerBase: *server.New(
			c.Timeouts.Prepare,
			func() error {
				if GRPC == nil {
					return trace.New(ErrNoGRPC)
				}
				return serveGRPC(
					GRPC,
					c,
					log,
					func() {
						log.Println("✗✗✗ Exit ✗✗✗")
					},
				)
			},
		),
	}

}

func serveGRPC(
	grpcServer *grpc.Server,
	c server.Configuration,
	log infrastructure.LoggerI,
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
