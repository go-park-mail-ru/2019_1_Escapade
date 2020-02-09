package models

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"

//easyjson:json
type LoadBalancer struct {
	ServiceName string `env:"name"`
	ServicePort int    `env:"port"`
	Entrypoint  string `env:"entrypoint"`
	Network     string `env:"network"`
}

func (l *LoadBalancer) Get() configuration.LoadBalancer {
	return configuration.LoadBalancer(*l)
}

func (l *LoadBalancer) Set(c configuration.LoadBalancer) {
	l.ServiceName = c.ServiceName
	l.ServicePort = c.ServicePort
	l.Entrypoint = c.Entrypoint
	l.Network = c.Network
}
