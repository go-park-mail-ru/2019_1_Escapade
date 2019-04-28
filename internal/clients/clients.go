package clients

import (
	session "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/proto"

	"google.golang.org/grpc"
)

type Clients struct {
	Session session.AuthCheckerClient
}

func Init(authConn *grpc.ClientConn) *Clients {
	return &Clients{Session: session.NewAuthCheckerClient(authConn)}
}
