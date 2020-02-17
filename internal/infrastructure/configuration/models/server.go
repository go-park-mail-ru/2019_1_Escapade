package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

//easyjson:json
type Server struct {
	Name           string   `json:"name" env:"name"`
	MaxConn        int      `json:"max_conn" env:"max_conn"`
	MaxHeaderBytes int      `json:"max_header_bytes"`
	Timeouts       Timeouts `json:"timeouts"`
	Port           int      `json:"port" env:"port"`
}

func (s *Server) Get() configuration.Server {
	return configuration.Server{
		Name:           s.Name,
		MaxConn:        s.MaxConn,
		MaxHeaderBytes: s.MaxHeaderBytes,
		Timeouts:       s.Timeouts.Get(),
		Port:           s.Port,
	}
}

func (s *Server) Set(c configuration.Server) {
	s.Name = c.Name
	s.MaxConn = c.MaxConn
	s.MaxHeaderBytes = c.MaxHeaderBytes
	s.Timeouts.Set(c.Timeouts)
	s.Port = c.Port

}

//easyjson:json
type Timeouts struct {
	//TTL domens.Duration `json:"ttl"`

	Read  models.Duration `json:"read"`
	Write models.Duration `json:"write"`
	Idle  models.Duration `json:"idle"`
	Wait  models.Duration `json:"wait"`
	Exec  models.Duration `json:"exec"`

	Prepare models.Duration `json:"prepare"`
}

func (t *Timeouts) Get() configuration.Timeouts {
	return configuration.Timeouts{
		//TTL: t.TTL.Duration,

		Read:  t.Read.Duration,
		Write: t.Write.Duration,
		Idle:  t.Idle.Duration,
		Wait:  t.Wait.Duration,
		Exec:  t.Exec.Duration,

		Prepare: t.Prepare.Duration,
	}
}

func (t *Timeouts) Set(c configuration.Timeouts) {
	//t.TTL.Duration = c.TTL

	t.Read.Duration = c.Read
	t.Write.Duration = c.Write
	t.Idle.Duration = c.Idle
	t.Wait.Duration = c.Wait
	t.Exec.Duration = c.Exec

	t.Prepare.Duration = c.Prepare

}
