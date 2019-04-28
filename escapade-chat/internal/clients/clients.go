package clients

import (
	session "escapade/internal/services/auth/proto"

	"google.golang.org/grpc"
)

type Clients struct {
	Session session.AuthCheckerClient
}

func Init(authConn *grpc.ClientConn) *Clients {
	return &Clients{Session: session.NewAuthCheckerClient(authConn)}
}
