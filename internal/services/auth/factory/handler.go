package factory

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/delivery/http/tokenhandler"
	ainfrastructure "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/infrastructure"
)

func NewHandler(
	r infrastructure.Router,
	mw []infrastructure.Middleware,
	server ainfrastructure.TokenServer,
	store ainfrastructure.TokenStore,
	trace infrastructure.ErrorTrace,
	logger infrastructure.Logger,
) http.Handler {
	r.PathHandler(("/swagger"), httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	var tokenHandler = tokenhandler.New(
		server,
		store,
		trace,
		logger,
	)

	var auth = r.PathSubrouter("/auth")

	auth.ANY("/login", tokenHandler.Login)
	auth.ANY("/auth", tokenHandler.Auth)
	auth.ANY("/delete", tokenHandler.Delete)
	auth.ANY("/test", tokenHandler.Test)
	auth.ANY("/token", tokenHandler.Create)
	auth.ANY("/authorize", tokenHandler.Authorize)

	auth.AddMiddleware(mw...)
	return r.Router()
}
