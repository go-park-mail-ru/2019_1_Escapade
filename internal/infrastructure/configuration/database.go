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
	ContextTimeout   time.Duration
}

type ConnectionStringRepository interface {
	Get() ConnectionString
	Set(ConnectionString)
}

type ConnectionString struct {
	User, Password, Database, Address string
}

func (str ConnectionString) ToString(driver string) string {
	return driver + "://" +
		str.User + ":" +
		str.Password + "@" + str.Address +
		"/" + str.Database + "?sslmode=disable"
}
