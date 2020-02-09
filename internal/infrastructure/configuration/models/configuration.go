package models

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"

type AllRepository interface {
	Get() All
	Set(All)
}

type All struct {
	Server           Server           `json:"server"`
	Auth             Auth             `json:"auth"`
	Database         Database         `json:"database"`
	Cors             Cors             `json:"cors"`
	LoadBalancer     LoadBalancer     `json:"loadbalancer"`
	Photo            Photo            `json:"photo"`
	ServiceDiscovery ServiceDiscovery `json:"service_discovery"`
}

func (all *All) Get() configuration.All {
	return configuration.All{
		Auth:             all.Auth.Get(),
		Cors:             all.Cors.Get(),
		Database:         all.Database.Get(),
		LoadBalancer:     all.LoadBalancer.Get(),
		Photo:            all.Photo.Get(),
		Server:           all.Server.Get(),
		ServiceDiscovery: all.ServiceDiscovery.Get(),
	}
}

func (all *All) Set(c configuration.All) {
	all.Auth.Set(c.Auth)
	all.Cors.Set(c.Cors)
	all.Database.Set(c.Database)
	all.LoadBalancer.Set(c.LoadBalancer)
	all.Photo.Set(c.Photo)
	all.Server.Set(c.Server)
	all.ServiceDiscovery.Set(c.ServiceDiscovery)
}
