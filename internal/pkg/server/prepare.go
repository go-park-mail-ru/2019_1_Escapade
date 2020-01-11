package server

import (
	"net/http"
	"time"
	"net"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

func ConfigureServer(handler http.Handler, serverConfig config.Server, port string) *http.Server {
	utils.Debug(false, "4. Configure server")

	// var execT = serverConfig.Timeouts.Exec.Duration

	// // в конфиг
	// var isHTTPS = false

	// if execT > time.Duration(time.Second) && !isHTTPS {
	// 	handler = http.TimeoutHandler(handler, execT, "ESCAPADE DEBUG Timeout!")
	// }

	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  serverConfig.Timeouts.Read.Duration,
		WriteTimeout: serverConfig.Timeouts.Write.Duration,
		IdleTimeout:  serverConfig.Timeouts.Idle.Duration,
		Handler:      handler,
		ConnState: func(c net.Conn, cs http.ConnState) {
			switch cs {
			case http.StateIdle, http.StateNew:
				c.SetReadDeadline(time.Now().Add(serverConfig.Timeouts.Idle.Duration))
			case http.StateActive:
				c.SetReadDeadline(time.Now().Add(serverConfig.Timeouts.Read.Duration))
			}
		},
		MaxHeaderBytes: serverConfig.MaxHeaderBytes,
	}
	return srv
}
