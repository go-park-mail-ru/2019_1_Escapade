package grpcclient

import (
	"errors"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/grpcclient/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
)

type GRPCCLient struct {
	SG               synced.SingleGoroutine
	Servers          []string
	GrcpConn         *grpc.ClientConn
	NameResolver     *testNameResolver
	ServiceDiscovery infrastructure.ServiceDiscovery
	Finish           chan interface{}

	C configuration.GRPCServer

	logger infrastructure.Logger

	errorCounterMutex *sync.Mutex
	_errorCounter     int // todo atomic?
}

func New(
	conf configuration.GRPCServerRepository,
	serviceDiscovery infrastructure.ServiceDiscovery,
	logger infrastructure.Logger,
) (*GRPCCLient, error) {
	// check configuration repository given
	if conf == nil {
		return nil, errors.New(ErrNoConfiguration)
	}
	c := conf.Get()

	// check service discovery given
	if serviceDiscovery == nil {
		return nil, errors.New(ErrNoServiceDiscovery)
	}

	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}

	var service = GRPCCLient{
		ServiceDiscovery:  serviceDiscovery,
		logger:            logger,
		C:                 c,
		NameResolver:      &testNameResolver{},
		errorCounterMutex: &sync.Mutex{},
	}

	var err error
	service.Servers, err = service.healthServers()
	if err != nil {
		return nil, err
	}
	service.Finish = make(chan interface{})
	service.SG = *synced.NewSingleGoroutine(c.Polling, service.poll)

	size := len(service.Servers)
	if size == 0 {
		logger.Println("cant get alive services")
		service.NameResolver.addr = DefaultAddress
	} else if size == 1 {
		service.NameResolver.addr = service.Servers[0]
	} else {
		service.append(service.Servers[1:size]...)
	}
	logger.Println("len(servers)", len(service.Servers))

	service.GrcpConn, err = grpc.Dial(
		service.NameResolver.addr,
		grpc.WithInsecure(),
		grpc.WithBalancer(grpc.RoundRobin(service.NameResolver)),
	)
	if err != nil {
		return nil, err
	}
	go service.runOnlineServiceDiscovery()
	return &service, nil
}

func (service *GRPCCLient) Close() error {
	err := service.GrcpConn.Close()
	service.SG.Close()
	service.Finish <- nil
	return err

}

func (service *GRPCCLient) runOnlineServiceDiscovery() {
	for {
		select {
		case <-service.SG.C():
			service.SG.Do()
		case <-service.Finish:
			close(service.Finish)
			return
		}
	}
}

func (service *GRPCCLient) healthServers() ([]string, error) {
	return service.ServiceDiscovery.Health(
		service.C.Name,
		service.C.Tag,
		PassingOnly,
	)
}

func (service *GRPCCLient) poll() {
	var (
		currAddrs    = service.generateMap(service.Servers)
		err          error
		errorDropped bool
	)
	service.Servers, err = service.healthServers()
	if err != nil {
		service.ErrorIncrese()
	}
	newAddrs := service.generateMap(service.Servers)
	if errorDropped = service.errorDrop(); errorDropped {
		newAddrs = make(map[string]struct{}, 0)
		service.Servers = []string{}
	}

	var updates []*naming.Update
	service.removeOld(currAddrs, newAddrs, updates)
	service.appendNew(currAddrs, newAddrs, updates)
	if len(updates) > 0 {
		service.NameResolver.w.inject(updates)
	}
	if errorDropped {
		service.poll()
	}
}

func (service *GRPCCLient) append(servers ...string) {
	var updates []*naming.Update
	for i := 1; i < len(servers); i++ {
		updates = append(updates, &naming.Update{
			Op:   naming.Add,
			Addr: servers[i],
		})
	}
	service.NameResolver.w.inject(updates)
}

func (service *GRPCCLient) removeOld(
	currAddrs, newAddrs map[string]struct{},
	updates []*naming.Update,
) {
	for addr := range currAddrs {
		if _, exist := newAddrs[addr]; !exist {
			updates = append(updates, &naming.Update{
				Op:   naming.Delete,
				Addr: addr,
			})
			delete(currAddrs, addr)
			service.logger.Println("remove", addr)
		}
	}
}

func (service *GRPCCLient) appendNew(currAddrs, newAddrs map[string]struct{}, updates []*naming.Update) {
	for addr := range newAddrs {
		if _, exist := currAddrs[addr]; !exist {
			updates = append(updates, &naming.Update{
				Op:   naming.Add,
				Addr: addr,
			})
			currAddrs[addr] = struct{}{}
			service.logger.Println("add", addr)
		}
	}
}

func (service *GRPCCLient) generateMap(arr []string) map[string]struct{} {
	arrInMap := make(map[string]struct{}, len(arr))
	for _, addr := range arr {
		arrInMap[addr] = struct{}{}
	}
	return arrInMap
}

func (service *GRPCCLient) ErrorIncrese() {
	service.errorCounterMutex.Lock()
	service._errorCounter++
	service.errorCounterMutex.Unlock()
}

func (service *GRPCCLient) errorDrop() bool {
	var v bool
	service.errorCounterMutex.Lock()
	if service._errorCounter >= service.C.CounterDrop {
		service._errorCounter = 0
		v = true
	}
	service.errorCounterMutex.Unlock()
	return v
}
