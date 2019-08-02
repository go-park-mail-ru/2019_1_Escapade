package clients

import (
	session "github.com/go-park-mail-ru/2019_1_Escapade/auth/server"
	pChat "github.com/go-park-mail-ru/2019_1_Escapade/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"

	"os"

	"google.golang.org/grpc"
)

type Clients struct {
	Session session.AuthCheckerClient
	Chat    pChat.ChatServiceClient
}

var ALL Clients

func Init(authConn *grpc.ClientConn) *Clients {
	return &Clients{Session: session.NewAuthCheckerClient(authConn)}
}

func (clients *Clients) InitChat(conn *grpc.ClientConn) {
	clients.Chat = pChat.NewChatServiceClient(conn)
}

func ServiceConnectionsInit(conf config.AuthClient) (authConn *grpc.ClientConn, err error) {

	authConn, err = grpc.Dial(
		os.Getenv(conf.URL),
		grpc.WithInsecure(),
	)
	if err != nil {
		return
	}

	//Other micro services conns wiil be here
	return
}
