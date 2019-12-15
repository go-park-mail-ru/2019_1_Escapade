package service

import (
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/handlers"
)

type Service struct {
	Chat     clients.ChatI
	Constant constants.RepositoryI
	Consul   server.ConsulServiceI
	Database *database.Input

	handler *handlers.GameHandler
}

func (s *Service) allSeT() bool {
	return s.Chat != nil && s.Constant != nil && s.Database != nil && s.Consul != nil
}

func (s *Service) Run(args *server.Args) error {
	if !s.allSeT() {
		return re.InterfaceIsNil()
	}

	c := args.Loader.Get()
	i := args.Input.GetData()
	if err := s.Chat.Init(s.Consul, c.Required); err != nil {
		return err
	}
	//defer chatService.Close()

	var gca = &handlers.ConfigurationArgs{
		C:         c,
		FieldPath: i.FieldPath,
		RoomPath:  i.RoomPath,
	}

	// start connection to database inside handlers
	s.handler = new(handlers.GameHandler)
	if err := s.handler.Init(s.Constant, s.Chat, gca, s.Database); err != nil {
		s.Chat.Close()
		return err
	}
	args.Handler = new(Router).Init(s.handler)
	//defer handler.Close()
	return nil
}

func (s *Service) Close() error {
	if !s.allSeT() {
		return re.InterfaceIsNil()
	}

	s.handler.Close()
	return s.Chat.Close()
}

//140 -> 64
