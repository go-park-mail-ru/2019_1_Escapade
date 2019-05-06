package clients

import (
	session "github.com/go-park-mail-ru/2019_1_Escapade/auth/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"

	"os"

	"google.golang.org/grpc"
)

type Clients struct {
	Session session.AuthCheckerClient
}

func Init(authConn *grpc.ClientConn) *Clients {
	return &Clients{Session: session.NewAuthCheckerClient(authConn)}
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
