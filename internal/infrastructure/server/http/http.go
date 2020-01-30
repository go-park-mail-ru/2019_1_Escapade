package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/netutil"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/server"
)

type ServerHTTP struct {
	server.ServerBase
}

func NewServerHTTP(
	c server.Configuration,
	handler http.Handler,
	log infrastructure.LoggerI,
	trace infrastructure.ErrorTrace,
) *ServerHTTP {
	if log == nil {
		log = &infrastructure.LoggerEmpty{}
	}
	if trace == nil {
		trace = &infrastructure.ErrorTraceDefault{}
	}
	var (
		server = &ServerHTTP{
			ServerBase: *server.New(
				c.Timeouts.Prepare,
				func() error {
					if handler == nil {
						return trace.New(ErrNoHandler)
					}
					return serveHTTP(
						configureServer(handler, c),
						c,
						log,
						func() {
							log.Println(false, "✗✗✗ Exit ✗✗✗")
						},
					)
				},
			),
		}
	)
	return server
}

func configureServer(
	handler http.Handler,
	c server.Configuration,
) *http.Server {
	// var execT = serverConfig.Timeouts.Exec.Duration

	// // в конфиг
	// var isHTTPS = false

	// if execT > time.Duration(time.Second) && !isHTTPS {
	// 	handler = http.TimeoutHandler(handler, execT, "ESCAPADE DEBUG Timeout!")
	// }

	return &http.Server{
		Addr:         c.Port,
		ReadTimeout:  c.Timeouts.Read,
		WriteTimeout: c.Timeouts.Write,
		IdleTimeout:  c.Timeouts.Idle,
		Handler:      handler,
		ConnState: func(conn net.Conn, cs http.ConnState) {
			switch cs {
			case http.StateIdle, http.StateNew:
				conn.SetReadDeadline(
					time.Now().Add(c.Timeouts.Idle),
				)
			case http.StateActive:
				conn.SetReadDeadline(
					time.Now().Add(c.Timeouts.Read),
				)
			}
		},
		MaxHeaderBytes: c.MaxHeaderBytes,
	}
}

func serveHTTP(
	srv *http.Server,
	c server.Configuration,
	log infrastructure.LoggerI,
	lastFunc func(),
) error {

	errChan := make(chan error)
	stopChan := make(chan os.Signal)
	defer func() {
		close(stopChan)
		close(errChan)
		lastFunc()
	}()

	signal.Notify(stopChan, os.Interrupt)

	l, err := net.Listen(Protocol, c.Port)
	if err != nil {
		log.Println("Listen error", err.Error())
		return err
	}

	defer l.Close()

	l = netutil.LimitListener(l, c.MaxConn)

	go func() {
		log.Println("✔✔✔ GO ✔✔✔")
		err := srv.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
			log.Println("Serving error:", err.Error())
		}
	}()
	waitTimeout := c.Timeouts.Wait
	ctx, cancel := context.WithTimeout(
		context.Background(),
		waitTimeout,
	)
	defer cancel()
	select {
	case err := <-errChan:
		log.Println("Fatal error: ", err.Error())
		return err
	case <-stopChan:
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Println("Shutdown error:", err.Error())
		}
	}
	<-ctx.Done()
	return nil
}
