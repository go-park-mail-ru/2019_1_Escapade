package database

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens"
)

type Configuration struct {
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type ConfigurationJSON struct {
	MaxOpenConns int             `json:"maxOpenConns"`
	MaxIdleConns int             `json:"maxIdleConns"`
	MaxLifetime  domens.Duration `json:"maxLifetime"`
}

func (c ConfigurationJSON) Get() Configuration {
	return Configuration{
		MaxOpenConns: c.MaxOpenConns,
		MaxIdleConns: c.MaxIdleConns,
		MaxLifetime:  c.MaxLifetime.Duration,
	}
}
