package service

import (
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/handlers"
)

type Service struct {
	Database *database.Input

	handler *handlers.Handlers
}

func (s *Service) allSeT() bool {
	return s.Database != nil
}

func (s *Service) Run(args *server.Args) error {
	if !s.allSeT() {
		return re.InterfaceIsNil()
	}

	s.handler = new(handlers.Handlers)
	err := s.handler.Init(args.Loader.Get(), s.Database)
	if err != nil {
		return err
	}

	args.Handler = new(Router).Init(s.handler)
	return nil
}

func (s *Service) Close() error {
	if !s.allSeT() {
		return re.InterfaceIsNil()
	}

	return s.handler.Close()
}
