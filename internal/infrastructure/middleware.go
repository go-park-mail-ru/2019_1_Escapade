package infrastructure

import "net/http"

type Middleware interface {
	Func(next http.Handler) http.Handler
}
