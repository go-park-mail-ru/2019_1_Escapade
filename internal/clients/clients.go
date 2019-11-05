package clients

import (
	//session "github.com/go-park-mail-ru/2019_1_Escapade/auth/server"

	pChat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"strconv"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"

	"google.golang.org/grpc"
)

type Clients struct {
	chatM  *sync.RWMutex
	_chat  pChat.ChatServiceClient
	consul *consulapi.Client
}

var ALL Clients

func (clients *Clients) Init( /*consulAddr string, ready chan error,
finish chan interface{}, conf config.Service*/) {

	clients.chatM = &sync.RWMutex{}
	/*
		for _, client := range conf.DependsOn {
			utils.Debug(false, "client name ", client)
			if client == "chat" {
				go clients.InitChat(client, consulAddr, ready, finish)
				<-ready
			}
		}*/
}

func (clients *Clients) AddChat(consulAddr string, finish chan interface{}) {
	ready := make(chan error)
	defer close(ready)

	go clients.InitChat("chat", consulAddr, ready, finish)
	<-ready
}

func (clients *Clients) InitChat(name string, consulAddr string, ready chan error, finish chan interface{}) {
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	consul, err := consulapi.NewClient(config)

	health, _, err := consul.Health().Service("chat", "", true, nil)
	if err != nil {
		utils.Debug(false, "cant get alive services")
	}

	servers := []string{}
	for _, item := range health {
		//item.Service.Address
		addr := "chat1" +
			":" + strconv.Itoa(item.Service.Port)
		servers = append(servers, addr)
	}

	nameResolver := &testNameResolver{}
	if len(servers) == 0 {
		utils.Debug(false, "cant get alive services")
		nameResolver.addr = ":0000"
	} else {
		nameResolver.addr = servers[0]
	}
	utils.Debug(false, "len(servers)", len(servers))

	grcpConn, err := grpc.Dial(
		nameResolver.addr,
		grpc.WithInsecure(),
		grpc.WithBalancer(grpc.RoundRobin(nameResolver)),
	)
	if err != nil {
		utils.Debug(false, "cant connect to grpc")
	}
	defer grcpConn.Close()

	if len(servers) > 1 {
		var updates []*naming.Update
		for i := 1; i < len(servers); i++ {
			updates = append(updates, &naming.Update{
				Op:   naming.Add,
				Addr: servers[i],
			})
		}
		nameResolver.w.inject(updates)
	}

	// тут мы будем периодически опрашивать консул на предмет изменений
	stopDiscovery := make(chan interface{}, 1)
	go runOnlineServiceDiscovery(nameResolver, consul, servers, stopDiscovery)

	if err != nil {
		utils.Debug(false, "Cant connect to chat service. Retry? ", err.Error())
		ready <- err
		return
	}
	defer grcpConn.Close()
	clients.setChat(grcpConn)
	ready <- err
	<-finish
	stopDiscovery <- nil
}

func (clients *Clients) setChat(conn *grpc.ClientConn) {
	clients.chatM.Lock()
	clients._chat = pChat.NewChatServiceClient(conn)
	clients.chatM.Unlock()
}

func (clients *Clients) Chat() pChat.ChatServiceClient {
	clients.chatM.RLock()
	v := clients._chat
	clients.chatM.RUnlock()
	return v
}

func runOnlineServiceDiscovery(nameResolver *testNameResolver,
	consul *consulapi.Client, servers []string, finish chan interface{}) {
	currAddrs := make(map[string]struct{}, len(servers))
	for _, addr := range servers {
		currAddrs[addr] = struct{}{}
	}
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			health, _, err := consul.Health().Service("chat", "", true, nil)
			if err != nil {
				utils.Debug(false, "cant get alive services")
			}

			newAddrs := make(map[string]struct{}, len(health))
			for _, item := range health {
				addr := item.Service.Address +
					":" + strconv.Itoa(item.Service.Port)
				newAddrs[addr] = struct{}{}
			}

			var updates []*naming.Update
			// проверяем что удалилось
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
			// проверяем что добавилось
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
			if len(updates) > 0 {
				nameResolver.w.inject(updates)
			}
		case <-finish:
			close(finish)
			return
		}
	}
}
