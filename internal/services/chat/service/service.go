package service

import (
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

type Service struct {
	Database *database.Input

	handler *handlers.Handler
}

func (s *Service) allSeT() bool {
	return s.Database != nil
}

func (s *Service) Run(args *server.Args) error {
	if !s.allSeT() {
		return re.InterfaceIsNil()
	}

	c := args.Loader.Get().DataBase
	s.handler = new(handlers.Handler)
	err := s.handler.Init(c, s.Database)
	if err != nil {
		return err
	}
	// defer handler.Close()

	args.GRPC = grpc.NewServer()
	proto.RegisterChatServiceServer(args.GRPC, s.handler)
	return nil
}

func (s *Service) Close() error {
	if !s.allSeT() {
		return re.InterfaceIsNil()
	}

	return s.handler.Close()
}

//140 -> 64
