package clients

import (
	//session "github.com/go-park-mail-ru/2019_1_Escapade/auth/server"
	pChat "github.com/go-park-mail-ru/2019_1_Escapade/chat/proto"
	config "github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"sync"

	"google.golang.org/grpc"
)

type Clients struct {
	//Session session.AuthCheckerClient

	chatM *sync.RWMutex
	_chat pChat.ChatServiceClient
}

var ALL Clients

func (clients *Clients) Init(ready chan error, finish chan interface{}, configClients ...config.Client) {
	clients.chatM = &sync.RWMutex{}
	utils.Debug(false, "init", len(configClients))
	for _, client := range configClients {
		utils.Debug(false, "client name ", client.Name)
		if client.Name == "chat" {
			go clients.InitChat(client.Address, ready, finish)
			<-ready
		}
	}
}

func (clients *Clients) InitChat(address string, ready chan error, finish chan interface{}) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		utils.Debug(false, "Cant connect to chat service. Retry? ", err.Error())
		ready <- err
		return
	}
	defer conn.Close()
	clients.setChat(conn)
	ready <- err
	<-finish
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

// func ServiceConnectionsInit(conf config.AuthClient) (authConn *grpc.ClientConn, err error) {

// 	authConn, err = grpc.Dial(
// 		os.Getenv(conf.URL),
// 		grpc.WithInsecure(),
// 	)
// 	if err != nil {
// 		return
// 	}

// 	//Other micro services conns wiil be here
// 	return
// }
