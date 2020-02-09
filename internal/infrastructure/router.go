package infrastructure

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

type Router interface {
	PathHandler(tpl string, handler http.Handler) Router
	PathHandlerFunc(
		tpl string,
		f func(http.ResponseWriter, *http.Request),
	) Router
	PathSubrouter(tpl string) Router
	GET(path string, f models.ResultFunc) Router
	POST(path string, f models.ResultFunc) Router
	PUT(path string, f models.ResultFunc) Router
	DELETE(path string, f models.ResultFunc) Router
	OPTIONS(path string, f models.ResultFunc) Router
	AddMiddleware(mwf ...Middleware) Router
	Router() http.Handler
}
