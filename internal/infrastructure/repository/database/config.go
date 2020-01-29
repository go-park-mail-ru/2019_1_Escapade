package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
)

type RepositoryConfig struct {
	entity.Database
}

func NewRepositoryConfig(c config.Configuration) *RepositoryConfig {
	var dc RepositoryConfig
	dc.MaxOpenConns = c.DataBase.MaxOpenConns
	dc.MaxIdleConns = c.DataBase.MaxIdleConns
	dc.MaxLifetime = c.DataBase.MaxLifetime.Duration
	return &dc
}

func (dc *RepositoryConfig) Get() entity.Database {
	return dc.Database
}
