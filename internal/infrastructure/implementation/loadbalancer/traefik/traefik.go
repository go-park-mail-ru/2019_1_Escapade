package traefik

import (
	"errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

type Traefik struct {
	configuration.LoadBalancer
}

// New instance of Traefik
func New(
	rep configuration.LoadBalancerRepository,
) (*Traefik, error) {
	if rep == nil {
		return nil, errors.New(ErrNoConfiguration)
	}
	return &Traefik{
		LoadBalancer: rep.Get(),
	}, nil
}

func (t *Traefik) RoutingTags() []string {
	var (
		name  = t.ServiceName
		port  = utils.String(t.ServicePort)
		entry = t.Entrypoint
		net   = t.Network
	)
	return []string{
		LabelEnable,
		LabelServices + name + ".loadbalancer.server.port=" + port,
		LabelRouters + name + ".service=" + name,
		LabelRouters + name + ".rule=PathPrefix(`/" + name + "`)",
		LabelRouters + name + ".entrypoints=" + entry,
		LabelNetwork + net,
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
