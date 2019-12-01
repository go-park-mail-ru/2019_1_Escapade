package server

import (
	"fmt"
	"os"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

//TODO
// realize docker, grpc, tcp, script check - https://www.consul.io/docs/agent/checks.html

/*
ConsulService responsible for registration, deregistration and
notification of service status(healthchecks) in the Consul

	ID - the ID of the service(get by calling ServiceID func) e.g api-14e7165cf399
	Name - the name of the service e.g api
	Address - the adress of the service(get by calling GetIP func)
	Port - port, lictened by the service
	Tags - consul tags as 'api', 'v2', 'traefic.enable=true' and so on
	TTL - interval of ttl sending to consul
	Check - the func, which return bool(is service working) and error
		based on the result of this function, the status of the service in consul

		true, nil - consulapi.HealthPassing
		false, nil - consulapi.HealthWarning
		*(any), error - consulapi.HealthCritical
	Checks - consul checks. Every instance of this type has TTL check.
		Also you can add http check if you call method .AddHTTPCheck
	ConsulAddr - address of consul. We need it to rejoin to consul client, if it fell
	initWeight - the initial weight of service for the load balancer
	_currentWeight - the current weight of service for the load balancer.  Protected by mutex!
	_client - client of Consul. Protected by mutex!
	enableTraefik - the flag responsible for the use of tags by the Traefik
*/
type ConsulService struct {
	ID         string
	Name       string
	Address    string
	Port       int
	Tags       []string
	TTL        time.Duration
	Check      func() (bool, error)
	Checks     consulapi.AgentServiceChecks
	ConsulAddr string
	initWeight int

	clientM *sync.RWMutex
	_client *consulapi.Client

	currentM       *sync.RWMutex
	_currentWeight int

	finish        chan interface{}
	enableTraefik bool
}

type ConsulInput struct {
	Name          string
	Port          int
	Tags          []string
	TTL           time.Duration
	MaxConn       int
	ConsulHost    string
	ConsulPort    string
	Check         func() (bool, error)
	EnableTraefik bool
}

// InitConsulService return instance of ConsulService
func InitConsulService(input *ConsulInput) *ConsulService {

	if input.EnableTraefik {
		input.Tags = append(input.Tags,
			"traefik.enable=true",
			"traefik.port=80",
			"traefik.docker.network=backend",
			"traefik.backend.loadbalancer=drr",
			"traefik.backend.maxconn.amount="+utils.String(input.MaxConn),
			"traefik.backend.maxconn.extractorfunc=client.ip")
	} else {
		input.Tags = append(input.Tags, "traefik.enable=false")
	}

	return generateService(input)
}

func generateService(input *ConsulInput) *ConsulService {
	var (
		id      = ServiceID(input.Name)
		weight  = CountWeight()
		address = GetIP()
	)

	checks := []*consulapi.AgentServiceCheck{
		&consulapi.AgentServiceCheck{
			CheckID:                        "service:" + id,
			TTL:                            input.TTL.String(),
			DeregisterCriticalServiceAfter: time.Minute.String(),
		}}
	return &ConsulService{
		ID:      id,
		Name:    input.Name,
		Address: address,
		Port:    input.Port,

		currentM:       &sync.RWMutex{},
		_currentWeight: weight,

		clientM: &sync.RWMutex{},
		_client: nil,

		initWeight:    weight,
		Tags:          input.Tags,
		TTL:           input.TTL,
		Check:         input.Check,
		Checks:        checks,
		ConsulAddr:    input.ConsulHost + input.ConsulPort,
		enableTraefik: input.EnableTraefik,
	}
}

// get the consul client
func (cs *ConsulService) connect() error {
	var (
		config = &consulapi.Config{
			Address:   cs.ConsulAddr,
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
func (cs *ConsulService) register(tags ...string) error {
	if cs.enableTraefik {
		tags = append(tags, "traefik.backend.weight="+utils.String(cs.weight()))
	}

	var (
		try = 3
		err error
	)
	for try >= 0 {
		try--
		client := cs.client()
		if client != nil {
			if err = client.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
				ID:      cs.ID,
				Name:    cs.Name,
				Port:    cs.Port,
				Address: cs.Address,
				Tags:    append(cs.Tags, tags...),
				Checks:  cs.Checks, //https://www.consul.io/docs/agent/checks.html
			}); err == nil {
				break
			}
		}
		if err = cs.connect(); err != nil {
			utils.Debug(false, "cant connect to consul", err.Error())
		}
	}
	return err
}

// Run the update goroutine. Dont forget to call .Close() to stop it
func (cs *ConsulService) Run() error {
	if err := cs.register(); err != nil {
		utils.Debug(false, "cant add service to consul", err)
		return err
	}

	cs.finish = make(chan interface{}, 1)
	go cs.updateTTL()
	return nil
}

// Close stop sending TTL goroutine and deregister service
func (cs *ConsulService) Close() error {
	cs.finish <- nil
	return cs.client().Agent().ServiceDeregister(cs.ID)
}

// updateTTL update TTl in consul. Called as goroutine. Will
// stop when the signal come in the channel 'finish'
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

// checkAndSetWeight updates the service weight in the load balancer
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
	return cs.register()
}

// Warn mark service status as Warning
// this will reduce the weight of the service twice
func (cs *ConsulService) Warn(note string) error {
	cs.checkAndSetWeight(cs.initWeight / 2)
	return cs.client().Agent().WarnTTL("service:"+cs.ID, note)
}

// AddHTTPCheck add http check to consul
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

// update - send service status to Consul
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
	err = cs.client().Agent().UpdateTTL("service:"+cs.ID, message, status)
	if err != nil {
		utils.Debug(false, "agent of", cs.ID, " UpdateTTL error:", err.Error())
		cs.register()
	}
}

// ServiceID return id of the service
func ServiceID(serviceName string) string {
	return serviceName + "-" + os.Getenv("HOSTNAME")
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

func (cs *ConsulService) Health() *consulapi.Health {
	return cs.client().Health()
}
