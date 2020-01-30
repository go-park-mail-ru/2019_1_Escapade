package load_balancer
/*
import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
	"os"
)

type RepositoryEnv struct {
	entity.LoadBalancer
}

func NewLoadBalancerEnv(name string, port int) *RepositoryEnv {
	var lbe RepositoryEnv
	lbe.Entrypoint = os.Getenv("entrypoint")
	lbe.Network = os.Getenv("network")
	lbe.ServiceName = name
	lbe.ServicePort = port
	return &lbe
}

func (lbe *RepositoryEnv) Get() entity.LoadBalancer {
	return lbe.LoadBalancer
}
*/