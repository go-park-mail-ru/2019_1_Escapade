package server

import (
	"os"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

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

func (ci *ConsulInput) Init(input InputI, loader ConfigutaionLoaderI) *ConsulInput {
	conf := loader.Get().Server
	ci.Name = conf.Name
	ci.Port = input.Port()
	ci.TTL = conf.Timeouts.TTL.Duration
	ci.MaxConn = conf.MaxConn
	ci.ConsulHost = os.Getenv("CONSUL_ADDRESS")
	ci.ConsulPort = ":8500"
	ci.EnableTraefik = conf.EnableTraefik

	entrypoint := "http"
	if os.Getenv("IS_HTTPS") != "" {
		entrypoint = "https"
	}
	ci.addTags(entrypoint)
	return ci
}

func (ci *ConsulInput) addTags(entrypoint string) {
	ci.Tags = []string{ci.Name,
		"traefik.frontend.rule=Host:" + ci.Name + ".consul.localhost",
		"traefik.frontend.entryPoints=" + entrypoint}
	ci.addTraefikTags()
}

func (ci *ConsulInput) addTraefikTags() {
	if ci.EnableTraefik {
		ci.Tags = append(ci.Tags,
			"traefik.enable=true",
			"traefik.port=80",
			"traefik.docker.network=backend",
			"traefik.backend.loadbalancer=drr",
			"traefik.backend.maxconn.amount="+utils.String(ci.MaxConn),
			"traefik.backend.maxconn.extractorfunc=client.ip")
	} else {
		ci.Tags = append(ci.Tags, "traefik.enable=false")
	}
}
