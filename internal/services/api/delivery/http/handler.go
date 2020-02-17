package http

import (
	"errors"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/handler"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handlers struct {
	Game    GameUseCase
	Session SessionUseCase
	Image   ImageUseCase
	User    UserUseCase
	Users   UsersUseCase
}

type UseCases struct {
	Image  api.ImageUseCaseI
	User   api.UserUseCaseI
	Record api.RecordUseCaseI
}

func NewHandler(
	logger infrastructure.Logger,
	errTrace infrastructure.ErrorTrace,
	r infrastructure.Router,
	h *Handlers,
	uc *UseCases,
	nonAuth []infrastructure.Middleware,
	withAuth []infrastructure.Middleware,
) (http.Handler, error) {
	if r == nil {
		return nil, errors.New(ErrNoRouter)
	}

	if h == nil {
		return nil, errors.New(ErrNoHandlers)
	}

	if uc == nil {
		return nil, errors.New(ErrNoUsecases)
	}

	var options = handler.New(logger, errTrace).OPTIONS()

	r.PathHandler("/swagger", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	http.Handle("/metrics", promhttp.Handler())
	//r.PathHandler("/metrics", promhttp.Handler()) // todo убрать дублирование

	var api = r.PathSubrouter("/api")
	var apiWithAuth = r.PathSubrouter("/api")

	apiWithAuth.POST("/game", h.Game.OfflineSave)
	api.OPTIONS("/game", options)

	apiWithAuth.DELETE("/session", h.Session.Logout)
	api.POST("/session", h.Session.Login)
	api.OPTIONS("/session", options)

	api.POST("/user", h.User.CreateUser)
	api.DELETE("/user", h.User.DeleteUser)
	apiWithAuth.PUT("/user", h.User.UpdateProfile)
	apiWithAuth.GET("/user", h.User.GetMyProfile)
	api.OPTIONS("/user", options)

	// delete "/avatar/{name}" path
	api.GET("/avatar/{name}", h.Image.GetImage)
	api.OPTIONS("/avatar/{name}", options)

	apiWithAuth.POST("/avatar", h.Image.PostImage)
	api.OPTIONS("/avatar", options)

	api.GET("/users/{id}", h.Users.GetOneUser)
	api.OPTIONS("/users/{id}", options)

	api.GET("/users/pages/page", h.Users.GetUsers)
	api.OPTIONS("/users/pages/page", options)

	api.GET("/users/pages/amount", h.Users.GetUsersPageAmount)
	api.OPTIONS("/users/pages/amount", options)

	var srv = server.ServerAddr{}
	api.PathHandlerFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("all ok " + srv.IP(nil)))
	})

	api.AddMiddleware(nonAuth...)
	apiWithAuth.AddMiddleware(withAuth...)

	return r.Router(), nil
}

// 128 -> 88
