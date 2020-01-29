package infrastructure

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
)

type AuthService interface {
	HashPassword(password string) string
	Check(
		rw http.ResponseWriter,
		r *http.Request,
	) (string, error)

	CreateToken(
		rw http.ResponseWriter,
		name, password string,
	) error

	DeleteToken(
		rw http.ResponseWriter,
		r *http.Request,
	) error
}

type AuthServiceRepositoryI interface {
	Get() entity.Auth
}
