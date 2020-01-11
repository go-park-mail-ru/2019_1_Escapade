package service

import (
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"net/http"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/handlers"
)

type Service struct {
	Database *database.Input
	handler *handlers.Handlers
}

func (s *Service) Init(db *database.Input) *Service {
	s.handler = new(handlers.Handlers)
	s.Database = db
	return s
}

func (s *Service) Run(args *server.Args) error {
	if err := re.NoNil(s, s.Database, s.handler); err != nil {
		return err
	}
	return s.handler.OpenDB(args.Loader.Get(), s.Database)
}

func (s *Service) Router() http.Handler {
	return s.handler.Router()
}

func (s *Service) Close() error {
	return s.handler.Close()
}
