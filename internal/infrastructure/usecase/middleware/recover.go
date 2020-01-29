package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

type Recover struct{}

func NewRecover() *Recover {
	return new(Recover)
}

func (mw *Recover) Func(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			defer utils.CatchPanic("middleware.go Recover()")
			next.ServeHTTP(rw, r)
		},
	)
}
