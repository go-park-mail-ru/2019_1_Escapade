package server

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens"
)

type Configuration struct {
	Name           string
	MaxConn        int
	MaxHeaderBytes int
	Timeouts       Timeouts
	EnableTraefik  bool
	Port           string
}

//easyjson:json
type ConfigurationJSON struct {
	Name           string       `json:"name"`
	MaxConn        int          `json:"maxConn"`
	MaxHeaderBytes int          `json:"maxHeaderBytes"`
	Timeouts       TimeoutsJSON `json:"timeouts"`
	EnableTraefik  bool         `json:"enableTraefik"`
	Port           string       `env:"port"`
}

func (s ConfigurationJSON) Get() Configuration {
	return Configuration{
		Name:           s.Name,
		MaxConn:        s.MaxConn,
		MaxHeaderBytes: s.MaxHeaderBytes,
		Timeouts:       s.Timeouts.Get(),
		EnableTraefik:  s.EnableTraefik,
		Port:           s.Port,
	}
}

// Timeouts of the connection to the server
type Timeouts struct {
	TTL time.Duration

	Read  time.Duration
	Write time.Duration
	Idle  time.Duration
	Wait  time.Duration
	Exec  time.Duration

	Prepare time.Duration
}

//easyjson:json
type TimeoutsJSON struct {
	TTL domens.Duration `json:"ttl"`

	Read  domens.Duration `json:"read"`
	Write domens.Duration `json:"write"`
	Idle  domens.Duration `json:"idle"`
	Wait  domens.Duration `json:"wait"`
	Exec  domens.Duration `json:"exec"`

	Prepare domens.Duration `json:"prepare"`
}

func (t TimeoutsJSON) Get() Timeouts {
	return Timeouts{
		TTL: t.TTL.Duration,

		Read:  t.Read.Duration,
		Write: t.Write.Duration,
		Idle:  t.Idle.Duration,
		Wait:  t.Wait.Duration,
		Exec:  t.Exec.Duration,

		Prepare: t.Prepare.Duration,
	}
}
