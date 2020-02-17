package http

import "net/http"

//go:generate $GOPATH/bin/mockery -name "UserRepositoryI|RecordRepositoryI|ImageRepositoryI"

type RepositoryI interface {
	GetUserID(req *http.Request) (int, error)
	GetPar(req *http.Request, par string) string
	GetDifficult(req *http.Request) string
	GetSort(req *http.Request) string
	GetName(req *http.Request) (string, error)
	GetPage(req *http.Request) string
	GetPerPage(req *http.Request) string
	GetNameAndPage(req *http.Request) (int, string, error)
}
