package load_balancer

type Interface interface {
	Init(Input *Input) Interface
	RoutingTags() []string
	WeightTags(name, value string) []string
}
