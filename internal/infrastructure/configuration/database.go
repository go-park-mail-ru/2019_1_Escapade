package configuration

import (
	"time"
)

// DatabaseRepository manage getting and setting the configuration of Database
type DatabaseRepository interface {
	Get() Database
	Set(Database)
}

// Database is configuration to initialize infrastructure.Database
type Database struct {
	MaxOpenConns     int
	MaxIdleConns     int
	MaxLifetime      time.Duration
	ConnectionString ConnectionString
}

type ConnectionString struct {
	User, Password, Database, Address string
}
