package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
)

type Recover struct {
	log infrastructure.LoggerI
}

func New(log infrastructure.LoggerI) *Recover {
	return &Recover{
		log: log,
	}
}

func (mw *Recover) Func(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					mw.log.Println("Panic recovered: ", r)
				}
			}()
			next.ServeHTTP(rw, r)
		},
	)
}
