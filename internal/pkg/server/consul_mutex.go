package server

import (
	consulapi "github.com/hashicorp/consul/api"
)

func (cs *ConsulService) weight() int {
	cs.currentM.RLock()
	cWeight := cs._currentWeight
	cs.currentM.RUnlock()
	return cWeight
}

func (cs *ConsulService) setWeight(weight int) {
	if weight < 0 {
		return
	}
	cs.currentM.Lock()
	cs._currentWeight = weight
	cs.currentM.Unlock()
}

func (cs *ConsulService) client() *consulapi.Client {
	cs.clientM.RLock()
	client := cs._client
	cs.clientM.RUnlock()
	return client
}

func (cs *ConsulService) setClient(client *consulapi.Client) {
	cs.clientM.Lock()
	cs._client = client
	cs.clientM.Unlock()
}
