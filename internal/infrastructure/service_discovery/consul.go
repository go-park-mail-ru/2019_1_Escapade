package service_discovery

import (
	"os"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

//TODO
// realize docker, grpc, tcp, script check - https://www.consul.io/docs/agent/checks.html

/*
ConsulService responsible for registration, deregistration and
notification of service status(healthchecks) in the Consul

	tags - consul tags as 'api', 'v2', 'traefic.enable=true' and so on
	TTL - interval of ttl sending to consul
	Check - the func, which return bool(is service working) and error
		based on the result of this function, the status of the service in consul

		true, nil - consulapi.HealthPassing
		false, nil - consulapi.HealthWarning
		*(any), error - consulapi.HealthCritical
	checks - consul checks. Every instance of this type has TTL check.
		Also you can add http check if you call method .AddHTTPCheck
	initWeight - the initial weight of service for the load balancer
	_currentWeight - the current weight of service for the load balancer.  Protected by mutex!
	_client - client of Consul. Protected by mutex!
*/
type Consul struct {
	Input        *infrastructure.ServiceDiscoveryData
	loadBalancer infrastructure.LoadBalancerI
	checks       consulapi.AgentServiceChecks

	initWeight int

	clientM *sync.RWMutex
	_client *consulapi.Client

	currentM       *sync.RWMutex
	_currentWeight int

	finish chan interface{}
}

// Init initialize ConsulService
func NewConsul(input *infrastructure.ServiceDiscoveryData) *Consul {
	var cs = new(Consul)
	cs.Input = input
	cs.checks = make([]*consulapi.AgentServiceCheck, 1)

	cs.Input.ID = ServiceID(cs.Input.ServiceName)
	cs.checks[0] = &consulapi.AgentServiceCheck{
		CheckID:                        cs.Input.ID,
		TTL:                            cs.Input.TTL.String(),
		DeregisterCriticalServiceAfter: time.Minute.String(), // todo в конфиг
	}

	var weight = CountWeight()
	cs.currentM = &sync.RWMutex{}
	cs._currentWeight = weight
	cs.initWeight = weight

	cs.clientM = &sync.RWMutex{}
	cs._client = nil

	return cs
}

// get the consul client
func (cs *Consul) connect() error {
	var (
		config = &consulapi.Config{
			Address:   cs.Input.DiscoveryAddr,
			Scheme:    "http",
			Transport: cleanhttp.DefaultPooledTransport(),
		}
		client, err = consulapi.NewClient(config)
	)
	if err == nil {
		cs.setClient(client)
	} else {
		cs.setClient(nil)
	}
	return err
}

// register our service in consul
// you can pass any number of tags to the function, which will
// be added to consul along with those that were specified when
// creating ConsulService(but these new tags will not be saved
// in ConsulService, only in Consul)
//
func (cs *Consul) register(tags ...string) error {
	var (
		client = cs.client()
		err    error
	)
	if client == nil {
		if err = cs.connect(); err != nil {
			utils.Debug(false, "cant connect to consul", err.Error())
			return err
		}
	}
	if cs.loadBalancer != nil {
		tags = append(tags,
			cs.loadBalancer.WeightTags(
				cs.Input.ID,
				utils.String(cs.weight()),
			)...,
		)
	}
	tags = append(cs.Input.Tags, tags...)
	err = cs.client().Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      cs.Input.ID,
		Name:    cs.Input.ServiceName,
		Port:    cs.Input.ServicePort,
		Address: cs.Input.ServiceHost,
		Tags:    append(cs.Input.Tags, tags...),
		Checks:  cs.checks, //https://www.consul.io/docs/agent/checks.html
	})
	if err != nil {
		if err = cs.connect(); err != nil {
			utils.Debug(false, "cant connect to consul", err.Error())
		}
	}
	return err
}

// Run the update goroutine. Dont forget to call .Close() to stop it
func (cs *Consul) Run() error {
	utils.Debug(false, "try register")
	if err := cs.register(); err != nil {
		utils.Debug(false, "cant add service to consul", err)
		return err
	}
	utils.Debug(false, "done")

	cs.finish = make(chan interface{}, 1)
	go cs.updateTTL()
	return nil
}

// Close stop sending TTL goroutine and deregister service
func (cs *Consul) Close() error {
	cs.finish <- nil
	return cs.client().Agent().ServiceDeregister(cs.Input.ID)
}

func (cs *Consul) AddLoadBalancer(lb infrastructure.LoadBalancerI) {
	cs.Input.Tags = append(cs.Input.Tags, lb.RoutingTags()...)
	cs.loadBalancer = lb
}

// updateTTL update TTl in consul. Called as goroutine. Will
// stop when the signal come in the channel 'finish'
func (cs *Consul) updateTTL() {
	var ttl = cs.Input.TTL
	if ttl.Seconds() > 5 {
		ttl = ttl - 5*time.Second // TODO убрать костыль
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

// checkAndSetWeight updates the service weight in the load balancer
func (cs *Consul) checkAndSetWeight(weight int) error {
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
	return cs.register()
}

// Warn mark service status as Warning
// this will reduce the weight of the service twice
func (cs *Consul) Warn(note string) error {
	cs.checkAndSetWeight(cs.initWeight / 2)
	return cs.client().Agent().WarnTTL("service:"+cs.Input.ID, note)
}

// HTTPCheck return http check to consul
func (cs *Consul) AddCheckHTTP(scheme, path, timeout, interval string) {
	address := scheme + "://" + cs.Input.ServiceHost + ":" + utils.String(cs.Input.ServicePort) + path
	cs.checks = append(cs.checks, &consulapi.AgentServiceCheck{
		CheckID:  "service:" + cs.Input.ID + ":http",
		Timeout:  timeout,
		Interval: interval,
		Method:   "GET",
		HTTP:     address,
	})
}

// update - send service status to Consul
func (cs *Consul) update() {
	var (
		isWarning bool
		err       error
		status    = consulapi.HealthPassing
		message   = "Alive and reachable"
	)
	if cs.Input.Check != nil {
		isWarning, err = cs.Input.Check()
	}
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
	client := cs.client()
	if client == nil {
		utils.Debug(false, "client == nil")
		err = cs.connect()
		if err != nil {
			utils.Debug(false, "cant connect", err.Error())
			return
		}
	}
	utils.Debug(false, "UpdateTTL")
	err = cs.client().Agent().UpdateTTL(cs.Input.ID, message, status)
	if err != nil {
		cs.client().Agent().ServiceDeregister(cs.Input.ID)
		utils.Debug(false, "agent of", cs.Input.ID, " UpdateTTL error:", err.Error())
		cs.register()
	}
}

// CountWeight return weight of the service taking into
// account its type recorded in the environment variables
func CountWeight() int {
	var weight = 6
	if os.Getenv("PRIMARY") != "" {
		weight = 12
	}
	if os.Getenv("SECONDARY") != "" {
		weight = 4
	}
	return weight
}

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

// ServiceID return id of the service
func ServiceID(serviceName string) string {
	return serviceName + "-" + os.Getenv("HOSTNAME")
}
