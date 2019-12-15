package service

import (
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	user_db "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/oauth"
	ery_db "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"
	pg "github.com/vgarvardt/go-oauth2-pg"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
)

type Service struct {
	Database    *database.Input
	RepositoryI clients.RepositoryI

	eryDB      *ery_db.DB
	tokenStore *pg.TokenStore
	userDB     user_db.UserUseCaseI
}

func (s *Service) allSeT() bool {
	return s.Database != nil && s.RepositoryI != nil
}

func (s *Service) Run(args *server.Args) error {
	if !s.allSeT() {
		return re.InterfaceIsNil()
	}

	var (
		err     error
		db      = args.Loader.Get().DataBase
		clients = s.RepositoryI.Get()
		manager *manage.Manager
	)

	s.eryDB, err = ery_db.Init("postgres://eryuser:nopassword@pg-ery:5432/erybase?sslmode=disable",
		db.MaxOpenConns, db.MaxIdleConns, db.MaxLifetime.Duration)
	if err != nil {
		return err
	}

	manager, s.tokenStore, err = oauth.Init(args.Loader.Get(), clients)
	if err != nil {
		re.Do(s.eryDB.Close, &err)
		return err
	}

	s.userDB = new(user_db.UserUseCase).Init(s.Database.User, s.Database.Record)

	if err = s.userDB.Open(db, s.Database.Database); err != nil {
		re.Do(s.eryDB.Close, &err)
		re.Do(s.tokenStore.Close, &err)
		return err
	}
	srv := oauth.Server(s.userDB, s.eryDB, manager)

	args.Handler = new(Router).Init(srv, s.tokenStore)
	return nil
}

func (s *Service) Close() error {
	if !s.allSeT() {
		return re.InterfaceIsNil()
	}
	return re.Close(s.eryDB, s.tokenStore)
}

func getClient() []*models.Client {
	return []*models.Client{&models.Client{
		ID:     "1",
		Secret: "1",
		Domain: "api.consul.localhost",
	}}
}
