package server

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
)

type RepositoryConfig struct {
	entity.Server
}

func NewServerConfig(c *config.Configuration, port string) *RepositoryConfig {
	var sc RepositoryConfig
	sc.Prepare = c.Server.Timeouts.Prepare.Duration
	sc.Port = port
	return &sc
}

func (sc *RepositoryConfig) Get() entity.Server {
	return sc.Server
}
