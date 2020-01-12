package service_discovery

import (
	"os"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server/load_balancer"
)

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
type Input struct {
	ID string

	ServiceName string
	ServicePort int
	ServiceHost string

	TTL     time.Duration
	MaxConn int

	DiscoveryAddr string
	tags          []string

	Check func() (bool, error)

	LoadBalancer load_balancer.Interface 
}

// Init initialize ConsulInput
func (ci *Input) Init(name string, port int, host string,
	ttl time.Duration, maxConn int, addr string,
	check func() (bool, error)) *Input {
	ci.ServiceName = name
	ci.ServicePort = port
	ci.ServiceHost = host // GetIP()
	ci.TTL = ttl
	ci.MaxConn = maxConn
	ci.DiscoveryAddr = addr
	ci.Check = check
	ci.tags = []string{name}
	return ci
}

func (ci *Input) AddLoadBalancer(lb load_balancer.Interface) {
	ci.tags = append(ci.tags, lb.RoutingTags()...)
	ci.LoadBalancer = lb
}

// ServiceID return id of the service
func ServiceID(serviceName string) string {
	return serviceName + "-" + os.Getenv("HOSTNAME")
}
