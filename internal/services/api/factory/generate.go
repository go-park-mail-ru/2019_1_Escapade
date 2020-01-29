package factory

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http"
	apiuc "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http/usecases"

	req "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http/repository"
	db "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/usecase/database"
)

func NewHandler(
	auth infrastructure.AuthService,
	photo infrastructure.PhotoServiceI,
	trace infrastructure.ErrorTrace,
	database infrastructure.DatabaseI,
	router infrastructure.RouterI,
	nonAuth []infrastructure.MiddlewareI,
	withAuth []infrastructure.MiddlewareI,
	timeout time.Duration, subnet string,
) http.Handler {

	usecases := NewUseCasesDB(database, trace, timeout)
	hlers := NewHandlers(
		auth,
		photo,
		*usecases,
		req.NewRequestMux(),
	)
	handler := handlers.NewHandler(
		router,
		hlers,
		usecases,
		subnet,
		nonAuth,
		withAuth,
	)
	return handler
}

func NewHandlers(
	auth infrastructure.AuthService,
	photo infrastructure.PhotoServiceI,
	hlrs handlers.UseCases,
	rep handlers.RepositoryI,
) *handlers.Handlers {
	return &handlers.Handlers{
		Game:    apiuc.NewGameHandler(hlrs.Record),
		Session: apiuc.NewSessionHandler(hlrs.User, auth),
		Image: apiuc.NewImageHandler(
			hlrs.Image,
			rep,
			photo,
		),
		User: apiuc.NewUserHandler(
			hlrs.User,
			hlrs.Record,
			rep,
			auth,
			photo,
		),
		Users: apiuc.NewUsersHandler(hlrs.User, rep, photo),
	}
}

func NewUseCasesDB(
	database infrastructure.DatabaseI,
	trace infrastructure.ErrorTrace,
	timeout time.Duration,
) *handlers.UseCases {
	return &handlers.UseCases{
		Image:  db.NewImage(database, timeout),
		User:   db.NewUser(database, trace, timeout),
		Record: db.NewRecord(database, timeout),
	}
}
