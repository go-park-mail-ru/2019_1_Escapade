package configuration

import (
	"time"
)

type ServerRepository interface {
	Get() Server
	Set(Server)
}

type Server struct {
	Name           string
	MaxConn        int
	MaxHeaderBytes int
	Timeouts       Timeouts
	Port           int
}

type TimeoutsRepository interface {
	Get() Timeouts
	Set(Timeouts)
}

// Timeouts of the connection to the server
type Timeouts struct {
	//TTL time.Duration

	Read  time.Duration
	Write time.Duration
	Idle  time.Duration
	Wait  time.Duration
	Exec  time.Duration

	Prepare time.Duration
}
