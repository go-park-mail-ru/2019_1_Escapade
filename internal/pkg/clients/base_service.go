package clients

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// RequiredService that is required for the correct working of this one
//easyjson:json
type RequiredService struct {
	Name        string          `json:"name"`
	Polling     models.Duration `json:"polling"`
	CounterDrop int             `json:"drop"`
	Tag         string          `json:"tag"`
}

type BaseService struct {
	SG               synced.SingleGoroutine
	Servers          []string
	GrcpConn         *grpc.ClientConn
	NameResolver     *testNameResolver
	ServiceDiscovery infrastructure.ServiceDiscovery
	Finish           chan interface{}

	C RequiredService

	errorCounterMutex *sync.Mutex
	_errorCounter     int
}

func (service *BaseService) Init(serviceDiscovery infrastructure.ServiceDiscovery,
	required RequiredService) error {

	service.Finish = make(chan interface{})
	service.ServiceDiscovery = serviceDiscovery
	service.C = required
	service.SG = synced.SingleGoroutine{}
	service.SG.Init(required.Polling.Duration, service.poll)

	service.Servers, _ = service.healthServers()
	service.NameResolver = &testNameResolver{}

	service.errorCounterMutex = &sync.Mutex{}

	size := len(service.Servers)
	if size == 0 {
		utils.Debug(false, "cant get alive services")
		service.NameResolver.addr = ":0000"
	} else if size == 1 {
		service.NameResolver.addr = service.Servers[0]
	} else {
		service.append(service.Servers[1:size]...)
	}
	utils.Debug(false, "len(servers)", len(service.Servers))

	var err error
	service.GrcpConn, err = grpc.Dial(
		service.NameResolver.addr,
		grpc.WithInsecure(),
		grpc.WithBalancer(grpc.RoundRobin(service.NameResolver)),
	)
	if err != nil {
		return err
	}
	go service.runOnlineServiceDiscovery()
	return nil
}

func (service *BaseService) Close() error {
	err := service.GrcpConn.Close()
	service.SG.Close()
	service.Finish <- nil
	return err

}

func (service *BaseService) runOnlineServiceDiscovery() {
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

func (service *BaseService) healthServers() ([]string, error) {
	return service.ServiceDiscovery.Health(service.C.Name,
		service.C.Tag, true)
}

func (service *BaseService) poll() {
	var (
		currAddrs    = generateMap(service.Servers)
		err          error
		errorDropped bool
	)
	service.Servers, err = service.healthServers()
	if err != nil {
		service.ErrorIncrese()
	}
	newAddrs := generateMap(service.Servers)
	if errorDropped = service.errorDrop(); errorDropped {
		newAddrs = make(map[string]struct{}, 0)
		service.Servers = []string{}
	}

	var updates []*naming.Update
	removeOld(currAddrs, newAddrs, updates)
	appendNew(currAddrs, newAddrs, updates)
	if len(updates) > 0 {
		service.NameResolver.w.inject(updates)
	}
	if errorDropped {
		service.poll()
	}
}

func (service *BaseService) append(servers ...string) {
	var updates []*naming.Update
	for i := 1; i < len(servers); i++ {
		updates = append(updates, &naming.Update{
			Op:   naming.Add,
			Addr: servers[i],
		})
	}
	service.NameResolver.w.inject(updates)
}

func removeOld(currAddrs, newAddrs map[string]struct{}, updates []*naming.Update) {
	for addr := range currAddrs {
		if _, exist := newAddrs[addr]; !exist {
			updates = append(updates, &naming.Update{
				Op:   naming.Delete,
				Addr: addr,
			})
			delete(currAddrs, addr)
			utils.Debug(false, "remove", addr)
		}
	}
}

func appendNew(currAddrs, newAddrs map[string]struct{}, updates []*naming.Update) {
	for addr := range newAddrs {
		if _, exist := currAddrs[addr]; !exist {
			updates = append(updates, &naming.Update{
				Op:   naming.Add,
				Addr: addr,
			})
			currAddrs[addr] = struct{}{}
			utils.Debug(false, "add", addr)
		}
	}
}

func generateMap(arr []string) map[string]struct{} {
	arrInMap := make(map[string]struct{}, len(arr))
	for _, addr := range arr {
		arrInMap[addr] = struct{}{}
	}
	return arrInMap
}

func (service *BaseService) ErrorIncrese() {
	service.errorCounterMutex.Lock()
	service._errorCounter++
	service.errorCounterMutex.Unlock()
}

func (service *BaseService) errorDrop() bool {
	var v bool
	service.errorCounterMutex.Lock()
	if service._errorCounter >= service.C.CounterDrop {
		service._errorCounter = 0
		v = true
	}
	service.errorCounterMutex.Unlock()
	return v
}
