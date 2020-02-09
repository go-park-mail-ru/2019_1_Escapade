package infrastructure

type LoadBalancer interface {
	RoutingTags() []string
	WeightTags(name, value string) []string
}

type LoadBalancerNil struct{}

func (*LoadBalancerNil) RoutingTags() []string {
	return []string{}
}
func (*LoadBalancerNil) WeightTags(
	name, value string,
) []string {
	return []string{}
}
