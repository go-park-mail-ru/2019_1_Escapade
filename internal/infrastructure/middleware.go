package infrastructure

import "net/http"

type MiddlewareI interface {
	Func(next http.Handler) http.Handler
}
