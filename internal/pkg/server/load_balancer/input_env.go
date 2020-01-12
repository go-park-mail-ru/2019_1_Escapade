package load_balancer

import (
	"os"
)

type InputEnv struct {
	Input
}

func (ie *InputEnv) Init(name string, port int) *Input {
	entrypoint := os.Getenv("entrypoint")
	network := os.Getenv("network")
	return ie.Input.Init(name, port, entrypoint, network)
}