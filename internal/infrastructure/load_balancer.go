package infrastructure

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"

type LoadBalancerI interface {
	RoutingTags() []string
	WeightTags(name, value string) []string
}

type LoadBalancerRepositoryI interface {
	Get() entity.LoadBalancer
}
