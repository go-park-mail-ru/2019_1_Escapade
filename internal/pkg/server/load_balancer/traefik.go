package load_balancer

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

type Traefik struct{
	Input *Input
}

func (t *Traefik) Init(input *Input) Interface {
	t.Input = input
	return t
}

func (t *Traefik) RoutingTags() []string {
	if t.Input == nil {
		return []string{}
	}
	var (
		name  = t.Input.ServiceName
		port  = utils.String(t.Input.ServicePort)
		entry = t.Input.Entrypoint
		net   = t.Input.Network
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

func (t *Traefik) WeightTags(serviceID, weight string) []string {
	prefix := "http.services."+t.Input.ServiceName+ ".weighted.services."
	return []string{
		prefix+"name="+serviceID,
		prefix+"weiht="+weight,
	}
}