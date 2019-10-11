package server

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/netutil"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func Server(r *mux.Router, serverConfig config.Server, isHTTP bool,
	port string) *http.Server {
	var (
		readTimeout  = time.Duration(serverConfig.ReadTimeoutS) * time.Second
		writeTimeout = time.Duration(serverConfig.WriteTimeoutS) * time.Second
		idleTimeout  = time.Duration(serverConfig.IdleTimeoutS) * time.Second
		execTimeout  = time.Duration(serverConfig.WaitTimeoutS) * time.Second
		handler      http.Handler
	)

	if serverConfig.WaitTimeoutS != 0 && isHTTP {
		handler = http.TimeoutHandler(r, execTimeout, "ESCAPADE DEBUG Timeout!")
	} else {
		handler = r
	}

	utils.Debug(false, "look", readTimeout, writeTimeout, idleTimeout, execTimeout)
	srv := &http.Server{
		Addr:           port,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
		Handler:        handler,
		MaxHeaderBytes: 1 << 15, // TODO в конфиг
	}
	return srv
}

func LaunchHTTP(server *http.Server, serverConfig config.Server, maxConn int,
	lastFunc func()) {
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
		return
	}

	defer l.Close()

	l = netutil.LimitListener(l, maxConn)

	go func() {
		utils.Debug(false, "✔✔✔ GO ✔✔✔")
		if err := server.Serve(l); err != nil && err != http.ErrServerClosed {
			errChan <- err
			utils.Debug(false, "Serving error:", err.Error())
		}
	}()
	waitTimeout := time.Duration(serverConfig.WaitTimeoutS) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	select {
	case err := <-errChan:
		utils.Debug(false, "Fatal error: ", err.Error())
		return
	case <-stopChan:
		err := server.Shutdown(ctx)
		if err != nil {
			utils.Debug(false, "Shutdown error:", err.Error())
		}
	}
	<-ctx.Done()
}

func LaunchGRPC(grpcServer *grpc.Server, port string, lastFunc func()) {
	errChan := make(chan error)
	stopChan := make(chan os.Signal)

	defer func() {
		close(stopChan)
		close(errChan)
		lastFunc()
	}()

	connectionCount := 20 // TODO в конфиг

	l, err := net.Listen("tcp", port)

	if err != nil {
		utils.Debug(true, "Listen error", err.Error())
		return
	}

	defer l.Close()

	l = netutil.LimitListener(l, connectionCount)

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
	case <-stopChan:
		grpcServer.GracefulStop()
	}
}
