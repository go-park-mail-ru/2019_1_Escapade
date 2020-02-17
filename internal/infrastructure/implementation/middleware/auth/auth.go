package miauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/handler"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
)

type Auth struct {
	handler.Handler

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
		Handler: *handler.New(logger, nil),
		auth:    auth,
		log:     logger,
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
				mw.SendResult(
					rw,
					mw.Fail(http.StatusUnauthorized, err),
				)
				return
			}
			ctx := context.WithValue(
				r.Context(),
				handler.ContextUserKey,
				userID,
			)

			mw.log.Println(
				"auth end",
				userID,
				handler.ContextUserKey,
			)
			next.ServeHTTP(rw, r.WithContext(ctx))
		},
	)
}
