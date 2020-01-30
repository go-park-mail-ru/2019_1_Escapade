package infrastructure

//go:generate $GOPATH/bin/mockery -name "Interface"

type ServiceDiscovery interface {
	Run() error
	Close() error
	AddLoadBalancer()

	Health(
		service,
		tag string,
		passingOnly bool,
	) ([]string, error)

	AddCheckHTTP(scheme, path, timeout, interval string)
}
