package load_balancer

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

type Traefik struct {
	entity.LoadBalancer
}

func New(
	rep infrastructure.LoadBalancerRepositoryI,
) *Traefik {
	return &Traefik{
		LoadBalancer: rep.Get(),
	}
}

func (t *Traefik) RoutingTags() []string {
	var (
		name  = t.ServiceName
		port  = utils.String(t.ServicePort)
		entry = t.Entrypoint
		net   = t.Network
	)
	return []string{
		"traefik.enable=true",
		"traefik.http.services." + name + ".loadbalancer.server.port=" + port,
		"traefik.http.routers." + name + ".service=" + name,
		"traefik.http.routers." + name + ".rule=PathPrefix(`/" + name + "`)",
		"traefik.http.routers." + name + ".entrypoints=" + entry,
		"traefik.docker.network=" + net,
	}
}

func (t *Traefik) WeightTags(
	serviceID, weight string,
) []string {
	prefix := "http.services." + t.ServiceName + ".weighted.services."
	return []string{
		prefix + "name=" + serviceID,
		prefix + "weiht=" + weight,
	}
}
