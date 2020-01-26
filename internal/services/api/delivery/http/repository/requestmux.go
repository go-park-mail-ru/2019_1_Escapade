package repository

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

type RequestMux struct{}

func NewRequestMux() *RequestMux {
	return &RequestMux{}
}

func (r *RequestMux) GetUserID(req *http.Request) (int, error) {
	return ih.IntFromPath(req, VarID, 1, re.ErrorInvalidUserID())
}

func (r *RequestMux) GetPar(req *http.Request, par string) string {
	return ih.StringFromPath(req, par, VarParDefault)
}

func (r *RequestMux) GetDifficult(req *http.Request) string {
	return ih.StringFromPath(req, VarDifficult, VarDifficultDefault)
}

func (r *RequestMux) GetSort(req *http.Request) string {
	return ih.StringFromPath(req, VarSort, VarSortDefault)
}

func (r *RequestMux) GetPage(req *http.Request) string {
	return ih.StringFromPath(req, VarPage, "0")
}

func (r *RequestMux) GetPerPage(req *http.Request) string {
	return ih.StringFromPath(req, VarPerPage, "")
}

func (r *RequestMux) GetName(req *http.Request) (username string, err error) {

	vars := mux.Vars(req)

	if username = vars[VarName]; username == "" {
		return "", re.ErrorInvalidName()
	}

	return
}

func (r *RequestMux) GetNameAndPage(req *http.Request) (page int, username string, err error) {
	vars := mux.Vars(req)

	if username = vars[VarName]; username == "" {
		return 0, "", re.ErrorInvalidName()
	}

	if vars[VarPage] == "" {
		page = 1
	} else {
		if page, err = strconv.Atoi(vars[VarPage]); err != nil {
			return 0, username, re.ErrorInvalidPage()
		}
		if page < 1 {
			page = 1
		}
	}
	return
}
