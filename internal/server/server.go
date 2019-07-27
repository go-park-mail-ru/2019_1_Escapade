package server

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
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

func InterruptHandlet(server *http.Server, serverConfig config.ServerConfig) {
	waitTimeout := time.Duration(serverConfig.WaitTimeoutS) * time.Second

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	go func() {
		err := server.Shutdown(ctx)
		if err != nil {
			utils.Debug(false, "Shutdown error:", err.Error())
		}
	}()
	<-ctx.Done()
	utils.Debug(false, "shutting down")
}
