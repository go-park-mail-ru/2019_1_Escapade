package entity

type LoadBalancer struct {
	ServiceName string
	ServicePort int
	Entrypoint  string
	Network     string
}
