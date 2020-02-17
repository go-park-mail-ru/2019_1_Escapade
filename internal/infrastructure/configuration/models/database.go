package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// Database representation of the configuration.Database as a json model
//easyjson:json
type Database struct {
	MaxOpenConns     int              `json:"max_open_conns"`
	MaxIdleConns     int              `json:"max_idle_conns"`
	MaxLifetime      models.Duration  `json:"max_life_time"`
	ConnectionString ConnectionString `json:"connection_string"`
	ContextTimeout   models.Duration  `json:"context_timeout"`
}

// Get configuration.Database from json model
// implementation of DatabaseRepository
func (d *Database) Get() configuration.Database {
	return configuration.Database{
		MaxOpenConns:     d.MaxOpenConns,
		MaxIdleConns:     d.MaxIdleConns,
		MaxLifetime:      d.MaxLifetime.Duration,
		ConnectionString: d.ConnectionString.Get(),
		ContextTimeout:   d.ContextTimeout.Duration,
	}
}

// Set data from configuration.Database
// implementation of DatabaseRepository
func (d *Database) Set(c configuration.Database) {
	d.MaxOpenConns = c.MaxOpenConns
	d.MaxIdleConns = c.MaxIdleConns
	d.MaxLifetime.Duration = c.MaxLifetime
	d.ConnectionString.Set(c.ConnectionString)
	d.ContextTimeout.Duration = c.ContextTimeout
}

//easyjson:json
type ConnectionString struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Address  string `json:"address" env:"db_address"`
}

func (cs *ConnectionString) Get() configuration.ConnectionString {
	return configuration.ConnectionString(*cs)
}

// Set data from configuration.Database
// implementation of DatabaseRepository
func (cs *ConnectionString) Set(c configuration.ConnectionString) {
	cs.User = c.User
	cs.Password = c.Password
	cs.Address = c.Address
	cs.Database = c.Database
}
