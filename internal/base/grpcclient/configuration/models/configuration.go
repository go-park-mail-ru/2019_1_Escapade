package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/grpcclient/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// RequiredService that is required for the correct working of this one
//easyjson:json
type GRPCServer struct {
	Name        string          `json:"name"`
	Polling     models.Duration `json:"polling"`
	CounterDrop int             `json:"drop"`
	Tag         string          `json:"tag"`
}

func (g *GRPCServer) Get() configuration.GRPCServer {
	return configuration.GRPCServer{
		Name:        g.Name,
		Polling:     g.Polling.Duration,
		CounterDrop: g.CounterDrop,
		Tag:         g.Tag,
	}
}

func (g *GRPCServer) Set(c configuration.GRPCServer) {
	g.Name = c.Name
	g.Polling.Duration = c.Polling
	g.CounterDrop = c.CounterDrop
	g.Tag = c.Tag

}
