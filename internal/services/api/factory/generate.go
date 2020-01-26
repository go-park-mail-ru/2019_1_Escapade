package factory

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/router"
	handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http"
	apiuc "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http/usecases"

	req "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http/repository"
	db "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/usecase/database"
)

func NewHandler(c *config.Configuration,
	dbI idb.Interface, timeout time.Duration,
	r router.Interface, subnet string) http.Handler {

	usecases := NewUseCasesDB(dbI, timeout)
	authToken := config.NewAuthToken(c.Auth, c.AuthClient, c.Cookie)
	hlrs := NewHandlers(*authToken, *usecases, req.NewRequestMux())
	handler := handlers.NewHandler(r, hlrs, usecases, authToken, c.Cors, subnet)
	return handler
}

func NewHandlers(cfg config.AuthToken,
	hlrs handlers.UseCases,
	rep handlers.RepositoryI) *handlers.Handlers {
	return &handlers.Handlers{
		Game:    apiuc.NewGameHandler(hlrs.Record),
		Session: apiuc.NewSessionHandler(cfg, hlrs.User),
		Image:   apiuc.NewImageHandler(hlrs.Image, rep),
		User:    apiuc.NewUserHandler(cfg, hlrs.User, hlrs.Record, rep),
		Users:   apiuc.NewUsersHandler(hlrs.User, rep),
	}
}

func NewUseCasesDB(dbI idb.Interface, timeout time.Duration) *handlers.UseCases {
	return &handlers.UseCases{
		Image:  db.NewImage(dbI, timeout),
		User:   db.NewUser(dbI, timeout),
		Record: db.NewRecord(dbI, timeout),
	}
}
