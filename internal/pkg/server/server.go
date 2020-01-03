package server

import (
	"context"
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// LaunchHTTP launch http server
func LaunchHTTP(server *http.Server, serverConfig config.Server,
	lastFunc func()) error {

	errChan := make(chan error)
	stopChan := make(chan os.Signal)
	defer func() {
		close(stopChan)
		close(errChan)
		lastFunc()
	}()

	signal.Notify(stopChan, os.Interrupt)

	l, err := net.Listen("tcp", server.Addr)
	if err != nil {
		utils.Debug(true, "Listen error", err.Error())
		return err
	}

	defer l.Close()

	l = netutil.LimitListener(l, serverConfig.MaxConn)

	go func() {
		utils.Debug(false, "✔✔✔ GO ✔✔✔")
		if err := server.Serve(l); err != nil && err != http.ErrServerClosed {
			errChan <- err
			utils.Debug(false, "Serving error:", err.Error())
		}
	}()
	waitTimeout := serverConfig.Timeouts.Wait.Duration
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	select {
	case err := <-errChan:
		utils.Debug(false, "Fatal error: ", err.Error())
		return err
	case <-stopChan:
		err := server.Shutdown(ctx)
		if err != nil {
			utils.Debug(false, "Shutdown error:", err.Error())
		}
	}
	<-ctx.Done()
	return nil
}

// LaunchGRPC launch grpc server
func LaunchGRPC(grpcServer *grpc.Server, serverConfig config.Server, port string, lastFunc func()) error {
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

func Port(port string) (string, int, error) {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return port, 0, err
	}
	return ":" + port, intPort, err
}

func GetIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err.Error()
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
