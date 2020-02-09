package mirecover

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
)

// Recover is implementation of Middleware interface(package infrastructure)
// recover panic
type Recover struct {
	log infrastructure.Logger
}

// New instance of Recover
func New(logger infrastructure.Logger) *Recover {
	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}
	return &Recover{
		log: logger,
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
