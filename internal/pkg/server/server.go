package server

import (
	"context"
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
			if (subnet==nil) {
				ips+=" " + ipsting
			} else if (strings.HasPrefix(ipsting, *subnet)) {
				return ipsting
			}
		}
	}
	if (subnet==nil) {
		return ips
	}
	return "error: no networks. Change subnet!"
}
