package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

type ServiceDiscovery struct {
	ID string `json:"id"`

	ServiceName string `json:"name" env:"name"`
	ServicePort int    `json:"port" env:"port"`
	ServiceHost string `json:"host"`

	Subnet string `env:"subnet"`

	TTL             models.Duration `json:"ttl"`
	CriticalTimeout models.Duration `json:"critical"`
	HTTPInterval    models.Duration `json:"http_interval"`
	HTTPTimeout     models.Duration `json:"http_timeout"`

	MaxConn int `json:"max_conn" env:"max_conn"`

	DiscoveryAddress string   `json:"discovery_address" env:"discovery_address"`
	Tags             []string `json:"tags" env:"tags"`
}

func (s *ServiceDiscovery) Get() configuration.ServiceDiscovery {
	s.ID = configuration.ServiceID(s.ServiceName)
	s.ServiceHost = configuration.GetIP(&s.Subnet)
	return configuration.ServiceDiscovery{
		ID: s.ID,

		ServiceName: s.ServiceName,
		ServicePort: s.ServicePort,
		ServiceHost: s.ServiceHost,

		Subnet: s.Subnet,

		TTL:             s.TTL.Duration,
		CriticalTimeout: s.CriticalTimeout.Duration,
		HTTPInterval:    s.HTTPInterval.Duration,
		HTTPTimeout:     s.HTTPTimeout.Duration,

		MaxConn: s.MaxConn,

		DiscoveryAddress: s.DiscoveryAddress,
		Tags:             s.Tags,
	}
}

func (s *ServiceDiscovery) Set(c configuration.ServiceDiscovery) {
	s.ID = c.ID

	s.ServiceName = c.ServiceName
	s.ServicePort = c.ServicePort
	s.ServiceHost = c.ServiceHost

	s.Subnet = c.Subnet

	s.TTL.Duration = c.TTL
	s.CriticalTimeout.Duration = c.CriticalTimeout
	s.HTTPInterval.Duration = c.HTTPInterval
	s.HTTPTimeout.Duration = c.HTTPTimeout

	s.MaxConn = c.MaxConn

	s.DiscoveryAddress = c.DiscoveryAddress
	s.Tags = c.Tags
}
