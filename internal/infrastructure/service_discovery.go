package infrastructure

import (
	"os"
	"time"
)

//go:generate $GOPATH/bin/mockery -name "Interface"

type ServiceDiscovery interface {
	Run() error
	Close() error
	AddLoadBalancer()

	AddCheckHTTP(scheme, path, timeout, interval string)
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
type ServiceDiscoveryData struct {
	ID string

	ServiceName string
	ServicePort int
	ServiceHost string

	TTL     time.Duration
	MaxConn int

	DiscoveryAddr string
	Tags          []string

	Check func() (bool, error)
}

func NewServiceDiscoveryData(name string, port int, host string,
	ttl time.Duration, maxConn int, addr string,
	check func() (bool, error)) *ServiceDiscoveryData {
	var ci = new(ServiceDiscoveryData)
	ci.ServiceName = name
	ci.ServicePort = port
	ci.ServiceHost = host // GetIP()
	ci.TTL = ttl
	ci.MaxConn = maxConn
	ci.DiscoveryAddr = addr
	ci.Check = check
	ci.Tags = []string{name}
	return ci
}

// ServiceID return id of the service
func ServiceID(serviceName string) string {
	return serviceName + "-" + os.Getenv("HOSTNAME")
}
