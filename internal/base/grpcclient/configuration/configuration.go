package configuration

import "time"

type GRPCServerRepository interface {
	Get() GRPCServer
	Set(GRPCServer)
}

// Configuration required for the correct working of this one
type GRPCServer struct {
	Name        string
	Polling     time.Duration
	CounterDrop int
	Tag         string
}
