package infrastructure

import "os"

type LoadBalancerI interface {
	RoutingTags() []string
	WeightTags(name, value string) []string
}

// data

type LoadBalancerData struct {
	ServiceName string
	ServicePort int
	Entrypoint  string
	Network     string
}

func NewLoadBalancerData(name string, port int, entrypoint, network string) *LoadBalancerData {
	var input = new(LoadBalancerData)
	input.ServiceName = name
	input.ServicePort = port
	input.Entrypoint = entrypoint
	input.Network = network
	return input
}

func NewLoadBalancerDataEnv(name string, port int) *LoadBalancerData {
	entrypoint := os.Getenv("entrypoint")
	network := os.Getenv("network")
	return NewLoadBalancerData(name, port, entrypoint, network)
}
