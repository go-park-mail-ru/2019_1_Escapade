package infrastructure

type LoadBalancerI interface {
	RoutingTags() []string
	WeightTags(name, value string) []string
}

type LoadBalancerEmpty struct{}

func (*LoadBalancerEmpty) RoutingTags() []string {
	return []string{}
}
func (*LoadBalancerEmpty) WeightTags(
	name, value string,
) []string {
	return []string{}
}
