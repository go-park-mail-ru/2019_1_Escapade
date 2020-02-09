package miauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
)

type Auth struct {
	auth infrastructure.AuthService
	log  infrastructure.Logger
}

func New(
	auth infrastructure.AuthService,
	logger infrastructure.Logger,
) (*Auth, error) {
	// check auth service given
	if auth == nil {
		return nil, errors.New(ErrNoAuthService)
	}
	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}
	return &Auth{
		auth: auth,
		log:  logger,
	}, nil
}

func (mw *Auth) Func(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			mw.log.Println("auth start")
			var (
				maxLimit = 1
				userID   string
				err      error
			)
			for i := 0; i < maxLimit; i++ { // TODO убрать костыль
				mw.log.Println("auth check start")
				userID, err = mw.auth.Check(rw, r)
				mw.log.Println("auth check end")
				if err == nil {
					mw.log.Println("no error auth")
					break
				} else {
					mw.log.Println("error auth", err.Error())
				}
			}
			if err != nil {
				handlers.SendResult(
					rw,
					handlers.NewResult(
						http.StatusUnauthorized,
						nil,
						err,
					),
					mw.log,
				)
				return
			}
			ctx := context.WithValue(
				r.Context(),
				handlers.ContextUserKey,
				userID,
			)

			mw.log.Println(
				"auth end",
				userID,
				handlers.ContextUserKey,
			)
			next.ServeHTTP(rw, r.WithContext(ctx))
		},
	)
}
