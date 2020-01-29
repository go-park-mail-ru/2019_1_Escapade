package server

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
)

// ServerHTTPConfig http server configuration
type ServerHTTPConfig struct {
	entity.ServerHTTP
}

// NewServerHTTPConfig create new instance of ServerHTTPConfig
func NewServerHTTPConfig(c *config.Configuration, port string) *ServerHTTPConfig {
	var sc ServerHTTPConfig
	sc.Server = NewServerConfig(c, port).Get()
	sc.Read = c.Server.Timeouts.Read.Duration
	sc.Write = c.Server.Timeouts.Write.Duration
	sc.Idle = c.Server.Timeouts.Idle.Duration
	sc.MaxHeaderBytes = c.Server.MaxHeaderBytes
	return &sc
}

// Get confiruration to initialize ServerI
// implements ServerHTTPRepositoryI
func (sc *ServerHTTPConfig) Get() entity.ServerHTTP {
	return sc.ServerHTTP
}
