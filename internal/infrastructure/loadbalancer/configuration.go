package loadbalancer

type Configuration struct {
	ServiceName string
	ServicePort int
	Entrypoint  string `env:"entrypoint"`
	Network     string `env:"network"`
}
