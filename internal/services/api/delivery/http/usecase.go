package api

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

//go:generate $GOPATH/bin/mockery -name "UserUseCaseI|RecordUseCaseI|ImageUseCaseI"

type GameUseCaseI interface {
	OfflineSave(rw http.ResponseWriter, r *http.Request) models.RequestResult
}

type SessionUseCaseI interface {
	Login(rw http.ResponseWriter, r *http.Request) models.RequestResult
	Logout(rw http.ResponseWriter, r *http.Request) models.RequestResult
}

type ImageUseCaseI interface {
	PostImage(rw http.ResponseWriter, r *http.Request) models.RequestResult
	GetImage(rw http.ResponseWriter, r *http.Request) models.RequestResult
}

type UserUseCaseI interface {
	CreateUser(rw http.ResponseWriter, r *http.Request) models.RequestResult
	GetMyProfile(rw http.ResponseWriter, r *http.Request) models.RequestResult
	DeleteUser(rw http.ResponseWriter, r *http.Request) models.RequestResult
	UpdateProfile(rw http.ResponseWriter, r *http.Request) models.RequestResult
}

type UsersUseCaseI interface {
	GetUsers(rw http.ResponseWriter, r *http.Request) models.RequestResult
	GetOneUser(rw http.ResponseWriter, r *http.Request) models.RequestResult
	GetUsersPageAmount(rw http.ResponseWriter, r *http.Request) models.RequestResult
}
