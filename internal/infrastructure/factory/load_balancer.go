package factory

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	rep "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/repository/load_balancer"
	uc "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/usecase/load_balancer"
)

func NewLoadBalancer(name string, port int) infrastructure.LoadBalancerI {
	var r = rep.NewLoadBalancerEnv(name, port)
	var pg = uc.NewTraefik(r)
	return pg
}
