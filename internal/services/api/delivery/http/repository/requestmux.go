package repository

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
)

type RequestMux struct {
	trace infrastructure.ErrorTrace
}

func NewRequestMux(trace infrastructure.ErrorTrace) *RequestMux {
	return &RequestMux{
		trace: trace,
	}
}

func (r *RequestMux) GetUserID(req *http.Request) (int, error) {
	return ih.IntFromPath(
		req,
		VarID,
		1,
		r.trace.New(ErrInvalidID),
	)
}

func (r *RequestMux) GetPar(req *http.Request, par string) string {
	return ih.StringFromPath(req, par, VarParDefault)
}

func (r *RequestMux) GetDifficult(req *http.Request) string {
	return ih.StringFromPath(
		req,
		VarDifficult,
		VarDifficultDefault,
	)
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

func (r *RequestMux) GetName(req *http.Request) (string, error) {
	var username string
	vars := mux.Vars(req)

	if username = vars[VarName]; username == "" {
		return "", r.trace.New(ErrInvalidName)
	}

	return username, nil
}

func (r *RequestMux) GetNameAndPage(req *http.Request) (int, string, error) {
	var (
		page     int
		username string
		err      error
	)
	vars := mux.Vars(req)

	if username = vars[VarName]; username == "" {
		return 0, "", r.trace.New(ErrInvalidName)
	}

	if vars[VarPage] == "" {
		page = 1
	} else {
		if page, err = strconv.Atoi(vars[VarPage]); err != nil {
			return 0, username, r.trace.New(ErrInvalidPage)
		}
		if page < 1 {
			page = 1
		}
	}
	return page, username, err
}
