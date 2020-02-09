package configuration

// LoadBalancerRepository manage getting and setting the configuration of LoadBalancer
type LoadBalancerRepository interface {
	Get() LoadBalancer
	Set(LoadBalancer)
}

// LoadBalancer is configuration to initialize infrastructure.LoadBalancer
type LoadBalancer struct {
	ServiceName string
	ServicePort int
	Entrypoint  string
	Network     string
}
