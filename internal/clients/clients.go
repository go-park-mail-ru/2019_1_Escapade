package clients

import (
	session "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/proto"

	"google.golang.org/grpc"
	"os"
)

type Clients struct {
	Session session.AuthCheckerClient
}

func Init(authConn *grpc.ClientConn) *Clients {
	return &Clients{Session: session.NewAuthCheckerClient(authConn)}
}

func ServiceConnectionsInit() (authConn *grpc.ClientConn, err error) {
	if os.Getenv("AUTHSERVICE_URL") == "" {
		os.Setenv("AUTHSERVICE_URL", "localhost:3333")
	}
	authConn, err = grpc.Dial(
		os.Getenv("AUTHSERVICE_URL"),
		grpc.WithInsecure(),
	)
	if err != nil {
		return
	}

	//Other micro services conns wiil be here

	return
}
