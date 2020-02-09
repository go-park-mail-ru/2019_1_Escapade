package configuration

type AllRepository interface {
	Get() All
	Set(All)
}

type All struct {
	Auth             Auth
	Cors             Cors
	Database         Database
	LoadBalancer     LoadBalancer
	Photo            Photo
	Server           Server
	ServiceDiscovery ServiceDiscovery
}
