package consul

const (
	Scheme = "http"

	ErrHealthNil       = "health is nil"
	ErrClientNil       = "consul client is nil"
	ErrNoConfiguration = "Configuration not given"

	CheckService         = "service:"
	CheckServiceProtocol = ":http"
	CheckMethod          = "GET"
	HealthMessage        = "Alive and reachable"

	EnvHost      = "HOSTNAME"
	EnvPrimary   = "PRIMARY"
	EnvSecondary = "SECONDARY"
)
