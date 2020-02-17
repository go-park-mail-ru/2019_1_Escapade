package configuration

import (
	"os"
	"time"
)

type ServiceDiscoveryRepository interface {
	Get() ServiceDiscovery
	Set(ServiceDiscovery)
}

/*
	ID - the ID of the service(get by calling ServiceID func) e.g api-14e7165cf399
	ServiceName - the name of the service e.g api
	ServiceHost - the host of the service(get by calling GetIP func)
	ServicePort - port, lictened by the service
	tags - service discovery tags as 'api', 'v2', 'traefic.enable=true' and so on
	TTL - interval of ttl sending to service discovery
	Check - the func, which return bool(is service working) and error
		based on the result of this function
	DiscoveryAddr - address of service discovery. We need it to rejoin to service
		discovery client after any failing
*/
// Input configuration of the service for its registration in service discovery
type ServiceDiscovery struct {
	ID string

	ServiceName string
	ServicePort int
	ServiceHost string

	Subnet string

	TTL             time.Duration
	CriticalTimeout time.Duration
	HTTPInterval    time.Duration
	HTTPTimeout     time.Duration

	MaxConn int

	DiscoveryAddress string
	Tags             []string
}

// func NewConfiguration(
// 	name string,
// 	port int,
// 	host string,
// 	ttl time.Duration,
// 	maxConn int,
// 	addr string,
// ) *Configuration {
// 	return &Configuration{
// 		ID: ServiceID(name),

// 		ServiceName: name,
// 		ServicePort: port,
// 		ServiceHost: host,

// 		TTL:     ttl,
// 		MaxConn: maxConn,

// 		DiscoveryAddr: addr,
// 		Tags:          []string{name},
// 	}
// }

// ServiceID return id of the service
func ServiceID(serviceName string) string {
	return serviceName + "-" + os.Getenv("HOSTNAME")
}
