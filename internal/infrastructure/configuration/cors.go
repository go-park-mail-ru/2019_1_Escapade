package configuration

// CorsRepository manage getting and setting the
//   configuration of Cors middleware
type CorsRepository interface {
	Get() Cors
	Set(Cors)
}

// Cors is configuration to initialize implementation of
//   infrastructure.Middleware Cors
type Cors struct {
	Origins     []string
	Headers     []string
	Methods     []string
	Credentials string
}
