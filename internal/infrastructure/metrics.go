package infrastructure

type MetricsI interface {
	HitsInc(ip, status, path, method string)
	UsersInc(ip, path, method string)
	UsersDec(ip, path, method string)
}
