package server

import (
	"os"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// ConsulInput configuration of the service for its registration in consul
type ConsulInput struct {
	Name          string
	Port          int
	Tags          []string
	TTL           time.Duration
	MaxConn       int
	ConsulHost    string
	ConsulPort    string
	Check         func() (bool, error)
	EnableTraefik bool
}

// Init initialize ConsulInput
func (ci *ConsulInput) Init(input InputI, loader ConfigutaionLoaderI) *ConsulInput {
	conf := loader.Get().Server
	ci.Name = conf.Name
	ci.Port = input.Port()
	ci.TTL = conf.Timeouts.TTL.Duration
	ci.MaxConn = conf.MaxConn
	ci.ConsulHost = os.Getenv("CONSUL_ADDRESS")
	ci.ConsulPort = ":8500"
	ci.EnableTraefik = conf.EnableTraefik
	ci.Check = func() (bool, error) {return false,nil}

	entrypoint := "http"
	if os.Getenv("IS_HTTPS") != "" {
		entrypoint = "https"
	}
	ci.Tags = []string{ci.Name}
	ci.addTraefikTags(entrypoint)
	return ci
}

// adds tags to interact with Traffic
func (ci *ConsulInput) addTraefikTags(entrypoint string) {
	if ci.EnableTraefik {
		ci.Tags = append(ci.Tags,
			"traefik.frontend.rule=PathPrefixStrip:/"+ci.Name,
			"traefik.frontend.entryPoints="+entrypoint,
			"traefik.enable=true",
			"traefik.port=3001",
			"traefik.docker.network=backend-overlay",
			"traefik.backend.loadbalancer=drr",
			"traefik.backend.maxconn.amount="+utils.String(ci.MaxConn),
			"traefik.backend.maxconn.extractorfunc=client.ip")
	} else {
		ci.Tags = append(ci.Tags, "traefik.enable=false")
	}
}
