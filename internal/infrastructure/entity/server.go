package entity

import "time"

type Server struct {
	Prepare time.Duration
	Port    string
}

type ServerGRPC struct {
	Server
}

type ServerHTTP struct {
	Server
	Read           time.Duration
	Write          time.Duration
	Idle           time.Duration
	MaxHeaderBytes int
}
