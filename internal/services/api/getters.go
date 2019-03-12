package api

import (
	"errors"
	"escapade/internal/misc"
	re "escapade/internal/return_errors"
	"net/http"
	"strconv"

	//"reflect"

	"github.com/gorilla/mux"
)

func (h *Handler) getPage(r *http.Request) (page int, err error) {

	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if vars["page"] == "" {
		page = 1
	} else {
		if page, err = strconv.Atoi(vars["page"]); err != nil {
			err = errors.New("Error page")
			return
		}
		if page < 1 {
			page = 1
		}

	}
	return
}

func (h *Handler) getNameAndPage(r *http.Request) (page int, username string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if username = vars["name"]; username == "" {
		err = re.ErrorInvalidName()
		return
	}

	if vars["page"] == "" {
		page = 1
	} else {
		if page, err = strconv.Atoi(vars["page"]); err != nil {
			err = re.ErrorInvalidPage()
			return
		}
		if page < 1 {
			page = 1
		}

	}
	return
}

func (h *Handler) getNameFromCookie(r *http.Request) (username string, err error) {
	var sessionID string

	if sessionID, err = misc.GetSessionCookie(r); err != nil {
		err = errors.New("Authorization required")
		return
	}

	if username, err = h.DB.GetNameBySessionID(sessionID); err != nil {
		return
	}

	return
}

func (h *Handler) getUserIDFromCookie(r *http.Request) (userID int, err error) {
	var sessionID string

	if sessionID, err = misc.GetSessionCookie(r); err != nil {
		err = errors.New("Authorization required")
		return
	}

	if userID, err = h.DB.GetUserIdBySessionID(sessionID); err != nil {
		return
	}

	return
}
