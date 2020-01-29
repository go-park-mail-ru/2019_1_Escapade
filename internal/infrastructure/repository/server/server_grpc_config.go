package server

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
)

type ServerGRPCConfig struct {
	entity.ServerGRPC
}

func NewServerGRPCConfig(c *config.Configuration, port string) *ServerGRPCConfig {
	var sc ServerGRPCConfig
	sc.Server = NewServerConfig(c, port).Get()
	return &sc
}

func (sc *ServerGRPCConfig) Get() entity.ServerGRPC {
	return sc.ServerGRPC
}
