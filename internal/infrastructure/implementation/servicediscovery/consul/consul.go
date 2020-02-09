package consul

import (
	"errors"
	"os"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
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
	c *configuration.ServiceDiscovery

	loadBalancer infrastructure.LoadBalancer
	log          infrastructure.Logger
	trace        infrastructure.ErrorTrace

	checks consulapi.AgentServiceChecks

	initWeight int

	clientM *sync.RWMutex
	_client *consulapi.Client

	currentM       *sync.RWMutex
	_currentWeight int

	Check func() (bool, error)

	finish chan struct{}
}

// Init initialize ConsulService
func New(
	conf configuration.ServiceDiscoveryRepository,
	check func() (bool, error),
	loadBalancer infrastructure.LoadBalancer,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
) (*Consul, error) {
	// check configuration repository given
	if conf == nil {
		return nil, errors.New(ErrNoConfiguration)
	}
	var c = conf.Get()

	//overriding the nil value of LoadBalancer
	if loadBalancer == nil {
		loadBalancer = new(infrastructure.LoadBalancerNil)
	}

	c.Tags = append(c.Tags, loadBalancer.RoutingTags()...)

	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}

	//overriding the nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}

	var weight = CountWeight()
	return &Consul{
		c: &c,

		loadBalancer: loadBalancer,
		log:          logger,
		trace:        trace,

		checks: []*consulapi.AgentServiceCheck{
			&consulapi.AgentServiceCheck{
				CheckID:                        c.ID,
				TTL:                            c.TTL.String(),
				DeregisterCriticalServiceAfter: c.CriticalTimeout.String(),
			}},
		initWeight: weight,
		clientM:    &sync.RWMutex{},

		currentM:       &sync.RWMutex{},
		_currentWeight: weight,

		Check: check,
	}, nil
}

// get the consul client
func (cs *Consul) connect() error {
	var (
		config = &consulapi.Config{
			Address:   cs.c.DiscoveryAddress,
			Scheme:    Scheme,
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
			cs.log.Println("cant connect to consul", err)
			return err
		}
	}
	tags = append(tags,
		cs.loadBalancer.WeightTags(
			cs.c.ID,
			utils.String(cs.weight()),
		)...,
	)

	tags = append(cs.c.Tags, tags...)
	err = cs.client().Agent().ServiceRegister(
		&consulapi.AgentServiceRegistration{
			ID:      cs.c.ID,
			Name:    cs.c.ServiceName,
			Port:    cs.c.ServicePort,
			Address: cs.c.ServiceHost,
			Tags:    append(cs.c.Tags, tags...),
			Checks:  cs.checks, //https://www.consul.io/docs/agent/checks.html
		},
	)
	if err != nil {
		if err = cs.connect(); err != nil {
			cs.log.Println("cant connect to consul", err)
		}
	}
	return err
}

// Run the update goroutine. Dont forget to call .Close() to stop it
func (cs *Consul) Run() error {
	cs.log.Println("try register")
	if err := cs.register(); err != nil {
		cs.log.Println("cant add service to consul", err)
		return err
	}
	cs.log.Println("done")

	cs.finish = make(chan struct{}, 1)
	go cs.updateTTL()
	return nil
}

func (cs *Consul) Health(
	service, tag string,
	passingOnly bool,
) ([]string, error) {
	client := cs.client()
	if client == nil {
		return []string{}, cs.trace.New(ErrClientNil)
	}
	health := client.Health()
	if health == nil {
		return []string{}, cs.trace.New(ErrHealthNil)
	}
	entries, _, err := health.Service(
		service,
		tag,
		passingOnly,
		nil,
	)
	if err != nil {
		return []string{}, err
	}
	var addresses = make([]string, len(entries))
	for i, item := range entries {
		addresses[i] = item.Service.Address //service + ":"+item.Service.Port
	}
	return addresses, nil
}

// Close stop sending TTL goroutine and deregister service
func (cs *Consul) Close() error {
	cs.finish <- struct{}{}
	return cs.client().Agent().ServiceDeregister(cs.c.ID)
}

// updateTTL update TTl in consul. Called as goroutine. Will
// stop when the signal come in the channel 'finish'
func (cs *Consul) updateTTL() {
	var ttl = cs.c.TTL
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
	return cs.client().Agent().WarnTTL(CheckService+cs.c.ID, note)
}

// HTTPCheck return http check to consul
func (cs *Consul) AddCheckHTTP(
	scheme, path, timeout, interval string,
) {
	address := scheme + "://" + cs.c.ServiceHost +
		":" + utils.String(cs.c.ServicePort) + path
	cs.checks = append(
		cs.checks, &consulapi.AgentServiceCheck{
			CheckID: CheckService +
				cs.c.ID + CheckServiceProtocol,
			Timeout:  timeout,
			Interval: interval,
			Method:   CheckMethod,
			HTTP:     address,
		},
	)
}

// update - send service status to Consul
func (cs *Consul) update() {
	var (
		isWarning bool
		err       error
		status    = consulapi.HealthPassing
		message   = HealthMessage
	)
	if cs.Check != nil {
		isWarning, err = cs.Check()
	}
	if err != nil {
		message = err.Error()
		if isWarning {
			status = consulapi.HealthWarning
			cs.log.Println(
				"healthcheck function warning:",
				message,
			)
			cs.checkAndSetWeight(cs.initWeight / 2)
		} else {
			status = consulapi.HealthCritical
			cs.log.Println(
				"healthcheck function error:",
				message,
			)
		}
	} else {
		cs.checkAndSetWeight(cs.initWeight)
	}
	client := cs.client()
	if client == nil {
		cs.log.Println("client == nil")
		err = cs.connect()
		if err != nil {
			cs.log.Println("cant connect", err.Error())
			return
		}
	}
	cs.log.Println("UpdateTTL")
	err = cs.client().Agent().UpdateTTL(
		cs.c.ID,
		message,
		status,
	)
	if err != nil {
		cs.client().Agent().ServiceDeregister(cs.c.ID)
		cs.log.Println(
			"agent of", cs.c.ID,
			" UpdateTTL error:", err)
		cs.register()
	}
}

// CountWeight return weight of the service taking into
// account its type recorded in the environment variables
func CountWeight() int {
	var weight = 6
	if os.Getenv(EnvPrimary) != "" {
		weight = 12
	}
	if os.Getenv(EnvSecondary) != "" {
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
	return serviceName + "-" + os.Getenv(EnvHost)
}
