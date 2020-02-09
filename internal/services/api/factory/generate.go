package factory

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http"
	apiuc "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http/usecase"

	req "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http/repository"
	db "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/usecase/database"
)

func NewHandler(
	auth infrastructure.AuthService,
	photo infrastructure.PhotoService,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
	database infrastructure.Database,
	router infrastructure.Router,
	nonAuth []infrastructure.Middleware,
	withAuth []infrastructure.Middleware,
	timeout time.Duration,
) http.Handler {
	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}

	//overriding the nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}

	//overriding the nil value of Photo Service
	if photo == nil {
		photo = new(infrastructure.PhotoServiceNil)
	}

	usecases := newUseCasesDB(database, trace, timeout)
	hlers := newHandlers(
		auth,
		photo,
		*usecases,
		req.NewRequestMux(trace),
		trace,
		logger,
	)
	handler := handlers.NewHandler(
		router,
		hlers,
		usecases,
		nonAuth,
		withAuth,
	)
	return handler
}

func newHandlers(
	auth infrastructure.AuthService,
	photo infrastructure.PhotoService,
	hlrs handlers.UseCases,
	rep handlers.RepositoryI,
	trace infrastructure.ErrorTrace,
	logger infrastructure.Logger,
) *handlers.Handlers {
	return &handlers.Handlers{
		Game: apiuc.NewGameHandler(hlrs.Record, trace),
		Session: apiuc.NewSessionHandler(
			hlrs.User,
			auth,
			trace,
			logger,
		),
		Image: apiuc.NewImageHandler(
			hlrs.Image,
			rep,
			photo,
			trace,
		),
		User: apiuc.NewUserHandler(
			hlrs.User,
			hlrs.Record,
			rep,
			auth,
			photo,
			trace,
			logger,
		),
		Users: apiuc.NewUsersHandler(
			hlrs.User,
			rep,
			photo,
			trace,
		),
	}
}

func newUseCasesDB(
	database infrastructure.Database,
	trace infrastructure.ErrorTrace,
	timeout time.Duration,
) *handlers.UseCases {
	return &handlers.UseCases{
		Image:  db.NewImage(database, timeout),
		User:   db.NewUser(database, trace, timeout),
		Record: db.NewRecord(database, timeout),
	}
}
