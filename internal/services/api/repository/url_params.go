package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

type UrlParam struct {
	r *http.Request
}

func (*UrlParam) UserID(r *http.Request) (int, error) {
	return ih.IntFromPath(r, "id", 1, re.ErrorInvalidUserID())
}

func (*UrlParam) Par(r *http.Request, par string) string {
	return ih.StringFromPath(r, par, "-")
}

func (*UrlParam) Difficult(r *http.Request) string {
	return ih.StringFromPath(r, "difficult", "0")
}

func (*UrlParam) Sort(r *http.Request) string {
	return ih.StringFromPath(r, "sort", "time")
}

func (*UrlParam) Name(r *http.Request) (username string, err error) {

	vars := mux.Vars(r)

	if username = vars["name"]; username == "" {
		return "", re.ErrorInvalidName()
	}

	return
}

func (*UrlParam) NameAndPage(r *http.Request) (page int, username string, err error) {
	vars := mux.Vars(r)

	if username = vars["name"]; username == "" {
		return 0, "", re.ErrorInvalidName()
	}

	if vars["page"] == "" {
		page = 1
	} else {
		if page, err = strconv.Atoi(vars["page"]); err != nil {
			return 0, username, re.ErrorInvalidPage()
		}
		if page < 1 {
			page = 1
		}
	}
	return
}
