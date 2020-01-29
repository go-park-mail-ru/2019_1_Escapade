package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/netutil"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

type ServerHTTP struct {
	entity.ServerHTTP
	infrastructure.ServerBase
}

func NewServerHTTP(rep infrastructure.ServerHTTPRepositoryI, handler http.Handler) *ServerHTTP {
	var server = new(ServerHTTP)
	server.ServerHTTP = rep.Get()
	var (
		run = func() error {
			var (
				exitFunc = func() { utils.Debug(false, "✗✗✗ Exit ✗✗✗") }
				srv      = server.configureServer(handler, c, port)
			)
			return server.serveHTTP(srv, exitFunc)
		}
	)
	server.SetRun(run)
	return server
}

func (server *ServerHTTP) configureServer(handler http.Handler) *http.Server {
	// var execT = serverConfig.Timeouts.Exec.Duration

	// // в конфиг
	// var isHTTPS = false

	// if execT > time.Duration(time.Second) && !isHTTPS {
	// 	handler = http.TimeoutHandler(handler, execT, "ESCAPADE DEBUG Timeout!")
	// }

	srv := &http.Server{
		Addr:         server.Port,
		ReadTimeout:  server.Read,
		WriteTimeout: server.Write,
		IdleTimeout:  server.Idle,
		Handler:      handler,
		ConnState: func(c net.Conn, cs http.ConnState) {
			switch cs {
			case http.StateIdle, http.StateNew:
				c.SetReadDeadline(time.Now().Add(server.Idle))
			case http.StateActive:
				c.SetReadDeadline(time.Now().Add(server.Read))
			}
		},
		MaxHeaderBytes: server.MaxHeaderBytes,
	}
	return srv
}

func (server *ServerHTTP) serveHTTP(srv *http.Server, lastFunc func()) error {

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
