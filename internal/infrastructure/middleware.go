package infrastructure

import "net/http"

type MiddlewareI interface {
	Func(http.Handler) http.Handler
}
