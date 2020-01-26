package api

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
)

//go:generate $GOPATH/bin/mockery -name "UserUseCaseI|RecordUseCaseI|ImageUseCaseI"

type GameUseCaseI interface {
	OfflineSave(rw http.ResponseWriter, r *http.Request) handlers.Result
}

type SessionUseCaseI interface {
	Login(rw http.ResponseWriter, r *http.Request) handlers.Result
	Logout(rw http.ResponseWriter, r *http.Request) handlers.Result
}

type ImageUseCaseI interface {
	PostImage(rw http.ResponseWriter, r *http.Request) handlers.Result
	GetImage(rw http.ResponseWriter, r *http.Request) handlers.Result
}

type UserUseCaseI interface {
	CreateUser(rw http.ResponseWriter, r *http.Request) handlers.Result
	GetMyProfile(rw http.ResponseWriter, r *http.Request) handlers.Result
	DeleteUser(rw http.ResponseWriter, r *http.Request) handlers.Result
	UpdateProfile(rw http.ResponseWriter, r *http.Request) handlers.Result
}

type UsersUseCaseI interface {
	GetUsers(rw http.ResponseWriter, r *http.Request) handlers.Result
	GetOneUser(rw http.ResponseWriter, r *http.Request) handlers.Result
	GetUsersPageAmount(rw http.ResponseWriter, r *http.Request) handlers.Result
}
