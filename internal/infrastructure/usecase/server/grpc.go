package server

import (
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"

	"golang.org/x/net/netutil"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

type ServerGRPC struct {
	entity.ServerGRPC
	infrastructure.ServerBase
}

func NewServerGRPC(GRPC *grpc.Server) *ServerGRPC {
	var (
		server = new(ServerGRPC)
		run    = func() error {
			var (
				exitFunc = func() { utils.Debug(false, "✗✗✗ Exit ✗✗✗") }
			)
			return serveGRPC(GRPC, c, port, exitFunc)
		}
	)
	server.Init(run, c)
	return server
}

func serveGRPC(grpcServer *grpc.Server, serverConfig config.Server, port string, lastFunc func()) error {
	errChan := make(chan error)
	stopChan := make(chan os.Signal)

	defer func() {
		close(stopChan)
		close(errChan)
		lastFunc()
	}()

	l, err := net.Listen("tcp", port)

	if err != nil {
		utils.Debug(true, "Listen error", err.Error())
		return err
	}

	defer l.Close()

	l = netutil.LimitListener(l, serverConfig.MaxConn)

	signal.Notify(stopChan, os.Interrupt)

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			errChan <- err
			close(errChan)
			utils.Debug(false, "Serving error:", err.Error())
		}
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		utils.Debug(false, "Fatal error: ", err.Error())
		return err
	case <-stopChan:
		grpcServer.GracefulStop()
	}
	return nil
}
