package service_discovery

import (
	consulapi "github.com/hashicorp/consul/api"
)

func (cs *Consul) weight() int {
	cs.currentM.RLock()
	cWeight := cs._currentWeight
	cs.currentM.RUnlock()
	return cWeight
}

func (cs *Consul) setWeight(weight int) {
	if weight < 0 {
		return
	}
	cs.currentM.Lock()
	cs._currentWeight = weight
	cs.currentM.Unlock()
}

func (cs *Consul) client() *consulapi.Client {
	cs.clientM.RLock()
	client := cs._client
	cs.clientM.RUnlock()
	return client
}

func (cs *Consul) setClient(client *consulapi.Client) {
	cs.clientM.Lock()
	cs._client = client
	cs.clientM.Unlock()
}
