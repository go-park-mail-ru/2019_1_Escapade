package server

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

type ConsulService struct {
	ID             string
	Name           string
	Address        string
	Port           int
	Tags           []string
	TTL            time.Duration
	Check          func() (bool, error)
	Client         *consulapi.Client
	Checks         consulapi.AgentServiceChecks
	initWeight     int
	_currentWeight int
	currentM       *sync.RWMutex
	finish         chan interface{}
}

//TODO
// realize docker, grpc, tcp, http, script check - https://www.consul.io/docs/agent/checks.html

type some *int
type somes *[]*some

func InitConsulService(id, name, address string, port int,
	tags []string, ttl time.Duration, maxConn int, primary bool,
	check func() (bool, error)) *ConsulService {
	var weight = 4
	if primary {
		weight = 12
	}

	checks := []*consulapi.AgentServiceCheck{
		&consulapi.AgentServiceCheck{
			CheckID:                        "service:" + id,
			TTL:                            ttl.String(),
			DeregisterCriticalServiceAfter: time.Minute.String(),
		}}

	return &ConsulService{
		ID:             id,
		Name:           name,
		Address:        address,
		Port:           port,
		_currentWeight: weight,
		currentM:       &sync.RWMutex{},
		initWeight:     weight,
		Tags: append(tags,
			"traefik.backend.loadbalancer=drr",
			"traefik.backend.maxconn.amount="+utils.String(maxConn),
			"traefik.backend.maxconn.extractorfunc=client.ip",
		),
		TTL:    ttl,
		Check:  check,
		Checks: checks,
	}
}

func (cs *ConsulService) InitAgent(consulAddr, consulPort string) error {
	var (
		config = &consulapi.Config{
			Address:   consulAddr + consulPort,
			Scheme:    "http",
			Transport: cleanhttp.DefaultPooledTransport(),
		}
		client, err = consulapi.NewClient(config)
	)
	cs.Client = client
	return err
}

func (cs *ConsulService) register(tags ...string) error {
	return cs.Client.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      cs.ID,
		Name:    cs.Name,
		Port:    cs.Port,
		Address: cs.Address,
		Tags:    append(cs.Tags, tags...),
		//https://www.consul.io/docs/agent/checks.html
		// Check: &consulapi.AgentServiceCheck{
		// 	TTL:                            cs.TTL.String(),
		// 	DeregisterCriticalServiceAfter: time.Minute.String(),
		// },
		Checks: cs.Checks,
	})
}

func (cs *ConsulService) Run() error {
	err := cs.register("traefik.backend.weight=" + utils.String(cs.Weight()))
	if err != nil {
		utils.Debug(false, "cant add service to consul", err)
		return err
	}

	cs.finish = make(chan interface{}, 1)
	go cs.updateTTL()
	return nil
}

func (cs *ConsulService) Close() error {
	cs.finish <- nil
	return cs.Client.Agent().ServiceDeregister(cs.ID)
}

func (cs *ConsulService) updateTTL() {
	var ttl = cs.TTL
	if cs.TTL.Seconds() > 5 {
		ttl = ttl - 5*time.Second
	}
	ticker := time.NewTicker(ttl)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cs.update()
		case <-cs.finish:
			close(cs.finish)
			return
		}
	}
}

func (cs *ConsulService) Weight() int {
	cs.currentM.RLock()
	cWeight := cs._currentWeight
	cs.currentM.RUnlock()
	return cWeight
}

func (cs *ConsulService) SetWeight(weight int) {
	if weight < 0 {
		return
	}
	cs.currentM.Lock()
	cs._currentWeight = weight
	cs.currentM.Unlock()
}

// return if changed
func (cs *ConsulService) checkAndSetWeight(weight int) error {
	var done bool
	cs.currentM.Lock()
	if cs._currentWeight != weight {
		done = true
		cs._currentWeight = weight
	}
	cs.currentM.Unlock()
	if !done {
		return nil
	}
	return cs.register("traefik.backend.weight=" + utils.String(weight))
}

func (cs *ConsulService) Warn(note string) error {
	cs.checkAndSetWeight(cs.initWeight / 2)
	return cs.Client.Agent().WarnTTL("service:"+cs.ID, note)
}

func (cs *ConsulService) AddHTTPCheck(scheme, path string) {
	address := scheme + "://" + cs.Address + ":" + utils.String(cs.Port) + path
	fmt.Println("toooook:", address)
	cs.Checks = append(cs.Checks, &consulapi.AgentServiceCheck{
		CheckID:  "service:" + cs.ID + ":http",
		Timeout:  "1s",
		Interval: "10s",
		Method:   "GET",
		HTTP:     address,
	})
}

func (cs *ConsulService) update() {
	var (
		isWarning, err = cs.Check()
		status         = consulapi.HealthPassing
		message        = "Alive and reachable"
	)
	if err != nil {
		message = err.Error()
		if isWarning {
			status = consulapi.HealthWarning
			utils.Debug(false, "healthcheck function warning:", message)
			cs.checkAndSetWeight(cs.initWeight / 2)
		} else {
			status = consulapi.HealthCritical
			utils.Debug(false, "healthcheck function error:", message)
		}
	} else {
		cs.checkAndSetWeight(cs.initWeight)
	}
	err = cs.Client.Agent().UpdateTTL("service:"+cs.ID, message, status)
	if err != nil {
		utils.Debug(false, "agent of", cs.ID, " UpdateTTL error:", err.Error())
	}
}

// ConsulClient register service and start healthchecking
func ConsulClient(serviceAddress, serviceName, host, serviceID string, portInt int, tags []string,
	consulPort string, ttl time.Duration, check func() (bool, error),
	finish chan interface{}) (*consulapi.Client, error) {
	var (
		config = &consulapi.Config{
			Address:   host + consulPort,
			Scheme:    "http",
			Transport: cleanhttp.DefaultPooledTransport(),
		}
		consul, err = consulapi.NewClient(config)
	)
	if err != nil {
		return consul, err
	}

	tags = append(tags, "traefik.backend.loadbalancer=drr",
		"traefik.backend.weight=10")
	agent := consul.Agent()
	err = registerService(agent, serviceID, serviceName, serviceAddress, portInt, tags, ttl)
	if err != nil {
		utils.Debug(false, "cant add service to consul", err)
		return consul, err
	}

	go updateTTL(agent, serviceID, ttl, check, finish)
	return consul, nil
}

func registerService(agent *consulapi.Agent, serviceID, serviceName,
	serviceAddress string, port int, tags []string, ttl time.Duration) error {
	return agent.ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Port:    port,
		Address: serviceAddress,
		Tags:    tags,
		// https://www.consul.io/docs/agent/checks.html
		Check: &consulapi.AgentServiceCheck{
			TTL:                            ttl.String(),
			DeregisterCriticalServiceAfter: time.Minute.String(),
		},
		Weights: &consulapi.AgentWeights{
			Passing: 100,
			Warning: 1,
		},
	})
}

// ServiceID return id of service
func ServiceID(serviceName string) string {
	return serviceName + "-" + os.Getenv("HOSTNAME")
}

func updateTTL(agent *consulapi.Agent, serviceID string, ttl time.Duration,
	check func() (bool, error), finish chan interface{}) {
	if ttl.Seconds() > 5 {
		ttl = ttl - 5*time.Second
	}
	ticker := time.NewTicker(ttl)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			update(agent, serviceID, check)
		case <-finish:
			close(finish)
			return
		}
	}
}

func update(agent *consulapi.Agent, serviceID string, check func() (bool, error)) {
	var (
		isWarning, err = check()
		status         = consulapi.HealthPassing
		message        = "Alive and reachable"
	)
	if err != nil {
		message = err.Error()
		if isWarning {
			status = consulapi.HealthWarning
			utils.Debug(false, "healthcheck function warning:", message)
		} else {
			status = consulapi.HealthCritical
			utils.Debug(false, "healthcheck function error:", message)
		}
	}
	err = agent.UpdateTTL("service:"+serviceID, message, status)
	if err != nil {
		utils.Debug(false, "agent of", serviceID, " UpdateTTL error:", err.Error())
	}
}
