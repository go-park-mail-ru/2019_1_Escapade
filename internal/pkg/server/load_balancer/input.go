package load_balancer

type Input struct {
	ServiceName string
	ServicePort int
	Entrypoint  string
	Network     string
}

func (input *Input) Init(name string, port int, entrypoint,
	network string) *Input {
	input.ServiceName = name
	input.ServicePort = port
	input.Entrypoint = entrypoint
	input.Network = network
	return input
}