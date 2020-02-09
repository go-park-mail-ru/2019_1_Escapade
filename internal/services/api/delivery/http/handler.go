package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handlers struct {
	Game    GameUseCaseI
	Session SessionUseCaseI
	Image   ImageUseCaseI
	User    UserUseCaseI
	Users   UsersUseCaseI
}

type UseCases struct {
	Image  api.ImageUseCaseI
	User   api.UserUseCaseI
	Record api.RecordUseCaseI
}

func NewHandler(
	r infrastructure.Router,
	h *Handlers,
	uc *UseCases,
	nonAuth []infrastructure.Middleware,
	withAuth []infrastructure.Middleware,
) http.Handler {

	r.PathHandler("/swagger", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	http.Handle("/metrics", promhttp.Handler())
	r.PathHandler("/metrics", promhttp.Handler())

	var api = r.PathSubrouter("/api")
	var apiWithAuth = r.PathSubrouter("/api")

	apiWithAuth.POST("/game", h.Game.OfflineSave)
	api.OPTIONS("/game", handlers.OPTIONS())

	apiWithAuth.DELETE("/session", h.Session.Logout)
	api.POST("/session", h.Session.Login)
	api.OPTIONS("/session", handlers.OPTIONS())

	api.POST("/user", h.User.CreateUser)
	api.DELETE("/user", h.User.DeleteUser)
	apiWithAuth.PUT("/user", h.User.UpdateProfile)
	apiWithAuth.GET("/user", h.User.GetMyProfile)
	api.OPTIONS("/user", handlers.OPTIONS())

	// delete "/avatar/{name}" path
	api.GET("/avatar/{name}", h.Image.GetImage)
	api.OPTIONS("/avatar/{name}", handlers.OPTIONS())

	apiWithAuth.POST("/avatar", h.Image.PostImage)
	api.OPTIONS("/avatar", handlers.OPTIONS())

	api.GET("/users/{id}", h.Users.GetOneUser)
	api.OPTIONS("/users/{id}", handlers.OPTIONS())

	api.GET("/users/pages/page", h.Users.GetUsers)
	api.OPTIONS("/users/pages/page", handlers.OPTIONS())

	api.GET("/users/pages/amount", h.Users.GetUsersPageAmount)
	api.OPTIONS("/users/pages/amount", handlers.OPTIONS())

	api.PathHandlerFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("all ok " + server.GetIP(nil)))
	})

	api.AddMiddleware(nonAuth...)
	apiWithAuth.AddMiddleware(withAuth...)

	return r.Router()
}

// 128 -> 88
