package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

type TokenHandler interface {
	Create(
		w http.ResponseWriter,
		r *http.Request,
	) models.RequestResult

	Delete(
		w http.ResponseWriter,
		r *http.Request,
	) models.RequestResult

	Authorize(
		w http.ResponseWriter,
		r *http.Request,
	) models.RequestResult

	Auth(
		w http.ResponseWriter,
		r *http.Request,
	) models.RequestResult

	Login(
		w http.ResponseWriter,
		r *http.Request,
	) models.RequestResult

	Test(
		w http.ResponseWriter,
		r *http.Request,
	) models.RequestResult
}
