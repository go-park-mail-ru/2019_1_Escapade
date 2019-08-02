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

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func Server(r *mux.Router, serverConfig config.ServerConfig, port string) *http.Server {
	var (
		readTimeout  = time.Duration(serverConfig.ReadTimeoutS) * time.Second
		writeTimeout = time.Duration(serverConfig.WriteTimeoutS) * time.Second
		idleTimeout  = time.Duration(serverConfig.IdleTimeoutS) * time.Second
		execTimeout  = time.Duration(serverConfig.WaitTimeoutS) * time.Second
		handler      http.Handler
	)

	if serverConfig.WaitTimeoutS == 0 {
		handler = http.TimeoutHandler(r, execTimeout, "Timeout!")
	} else {
		handler = r
	}

	utils.Debug(false, "look", readTimeout, writeTimeout, idleTimeout, execTimeout)
	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      handler,
	}
	return srv
}

func LaunchHTTP(server *http.Server, serverConfig config.ServerConfig, lastFunc func()) {
	errChan := make(chan error)
	stopChan := make(chan os.Signal)

	signal.Notify(stopChan, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			errChan <- err
			utils.Debug(false, "Serving error:", err.Error())
		}
	}()

	defer func() {
		close(errChan)
		close(stopChan)
		lastFunc()
	}()

	waitTimeout := time.Duration(serverConfig.WaitTimeoutS) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)

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

	defer cancel()
	<-ctx.Done()

	go func() {
		err := server.Shutdown(ctx)
		if err != nil {
			utils.Debug(false, "Shutdown error:", err.Error())
		}
		lastFunc()
	}()
}

func LaunchGRPC(grpcServer *grpc.Server, lis net.Listener, lastFunc func()) {
	errChan := make(chan error)
	stopChan := make(chan os.Signal)

	signal.Notify(stopChan, os.Interrupt)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- err
			utils.Debug(false, "Serving error:", err.Error())
		}
	}()

	defer func() {
		grpcServer.GracefulStop()
		close(errChan)
		close(stopChan)
		lastFunc()
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		utils.Debug(false, "Fatal error: ", err.Error())
	case <-stopChan:
	}
}
