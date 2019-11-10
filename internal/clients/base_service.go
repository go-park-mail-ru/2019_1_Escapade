package clients

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/synced"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
)

type BaseService struct {
	sg           synced.SingleGoroutine
	servers      []string
	grcpConn     *grpc.ClientConn
	nameResolver *testNameResolver
	consul       *server.ConsulService
	finish       chan interface{}

	c config.RequiredService

	errorCounterMutex *sync.Mutex
	_errorCounter     int
}

func (service *BaseService) Init(consul *server.ConsulService,
	required config.RequiredService) error {

	service.finish = make(chan interface{})
	service.consul = consul
	service.c = required
	service.sg = synced.SingleGoroutine{}
	service.sg.Init(required.Polling.Duration, service.poll)

	service.servers, _ = service.healthServers()
	service.nameResolver = &testNameResolver{}

	service.errorCounterMutex = &sync.Mutex{}

	size := len(service.servers)
	if size == 0 {
		utils.Debug(false, "cant get alive services")
		service.nameResolver.addr = ":0000"
	} else if size == 1 {
		service.nameResolver.addr = service.servers[0]
	} else {
		service.append(service.servers[1:size]...)
	}
	utils.Debug(false, "len(servers)", len(service.servers))

	var err error
	service.grcpConn, err = grpc.Dial(
		service.nameResolver.addr,
		grpc.WithInsecure(),
		grpc.WithBalancer(grpc.RoundRobin(service.nameResolver)),
	)
	if err != nil {
		return err
	}
	go service.runOnlineServiceDiscovery()
	return nil
}

func (service *BaseService) Close() {
	service.grcpConn.Close()
	service.sg.Close()
	service.finish <- nil

}

func (service *BaseService) runOnlineServiceDiscovery() {
	for {
		select {
		case <-service.sg.C():
			service.sg.Do()
		case <-service.finish:
			close(service.finish)
			return
		}
	}
}

func (service *BaseService) healthServers() ([]string, error) {
	if service.consul.Health() == nil {
		return []string{}, re.ErrorServer()
	}
	health, _, err := service.consul.Health().Service(service.c.Name,
		service.c.Tag, true, nil)
	if err != nil {
		utils.Debug(false, "consul error", err.Error())
		return []string{}, err
	}
	servers := []string{}
	for _, item := range health {
		//item.Service.Address
		addr := service.c.Name + ":" + strconv.Itoa(item.Service.Port)
		fmt.Println("chat addr:", addr)
		servers = append(servers, addr)
	}
	return servers, nil
}

func (service *BaseService) poll() {
	var (
		currAddrs    = generateMap(service.servers)
		err          error
		errorDropped bool
	)
	service.servers, err = service.healthServers()
	if err != nil {
		service.errorIncrese()
	}
	newAddrs := generateMap(service.servers)
	if errorDropped = service.errorDrop(); errorDropped {
		newAddrs = make(map[string]struct{}, 0)
		service.servers = []string{}
	}

	var updates []*naming.Update
	removeOld(currAddrs, newAddrs, updates)
	appendNew(currAddrs, newAddrs, updates)
	if len(updates) > 0 {
		service.nameResolver.w.inject(updates)
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
	service.nameResolver.w.inject(updates)
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

func (service *BaseService) errorIncrese() {
	service.errorCounterMutex.Lock()
	service._errorCounter++
	service.errorCounterMutex.Unlock()
}

func (service *BaseService) errorDrop() bool {
	var v bool
	service.errorCounterMutex.Lock()
	if service._errorCounter >= service.c.CounterDrop {
		service._errorCounter = 0
		v = true
	}
	service.errorCounterMutex.Unlock()
	return v
}
