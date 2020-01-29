package infrastructure

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
)

type RouterI interface {
	PathHandler(tpl string, handler http.Handler) RouterI
	PathHandlerFunc(
		tpl string,
		f func(http.ResponseWriter, *http.Request),
	) RouterI
	PathSubrouter(tpl string) RouterI
	GET(path string, f handlers.ResultFunc) RouterI
	POST(path string, f handlers.ResultFunc) RouterI
	PUT(path string, f handlers.ResultFunc) RouterI
	DELETE(path string, f handlers.ResultFunc) RouterI
	OPTIONS(path string, f handlers.ResultFunc) RouterI
	AddMiddleware(mwf ...MiddlewareI) RouterI
	Router() http.Handler
}
