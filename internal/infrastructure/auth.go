package infrastructure

import "net/http"

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