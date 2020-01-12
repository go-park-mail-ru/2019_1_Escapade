package service_discovery

import (
	consulapi "github.com/hashicorp/consul/api"
)

//go:generate $GOPATH/bin/mockery -name "Interface"

type Interface interface {
	Init(input *Input) Interface
	Health() *consulapi.Health

	Data() *Input

	Run() error
	Close() error
}
