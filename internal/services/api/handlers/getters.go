package handlers

import (
	"net/http"
	"strconv"

	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"github.com/gorilla/mux"
)

func getUserID(r *http.Request) (int, error) {
	return ih.IntFromPath(r, "id", 1, re.ErrorInvalidUserID())
}

func getPar(r *http.Request, par string) string {
	return ih.StringFromPath(r, par, "-")
}

func getDifficult(r *http.Request) string {
	return ih.StringFromPath(r, "difficult", "0")
}

func getSort(r *http.Request) string {
	return ih.StringFromPath(r, "sort", "time")
}

func getName(r *http.Request) (username string, err error) {

	vars := mux.Vars(r)

	if username = vars["name"]; username == "" {
		return "", re.ErrorInvalidName()
	}

	return
}

func getNameAndPage(r *http.Request) (page int, username string, err error) {
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
