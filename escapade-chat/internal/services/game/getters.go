package game

import (
	"context"
	"escapade/internal/config"
	"escapade/internal/cookie"
	re "escapade/internal/return_errors"
	"net/http"
	"strconv"

	session "escapade/internal/services/auth/proto"
)

func getStringFromPath(r *http.Request, name string, defaultValue string) (str string) {
	str = defaultValue

	vals := r.URL.Query()
	keys, ok := vals[name]
	if ok {
		if len(keys) >= 1 {
			str = keys[0]
		}
	}
	return
}

func getIntFromPath(r *http.Request, name string,
	defaultVelue int, expected error) (val int, err error) {
	var str string
	if str = getStringFromPath(r, name, ""); str == "" {
		err = expected
		return
	}
	val = defaultVelue

	if val, err = strconv.Atoi(str); err != nil {
		err = expected
		return
	}
	if val < 0 {
		err = expected
		return
	}
	return
}

func (h *Handler) getUserID(r *http.Request) (id int, err error) {

	id, err = getIntFromPath(r, "id", 1, re.ErrorInvalidUserID())
	return
}

func (h *Handler) getPage(r *http.Request) (page int) {

	page, _ = getIntFromPath(r, "page", 1, nil)
	return
}

func (h *Handler) getUserIDFromCookie(r *http.Request, cc config.CookieConfig) (userID int, name string, err error) {
	sessionID, _ := cookie.GetSessionCookie(r, cc)

	ctx := context.Background()
	sess, err := h.Clients.Session.Check(ctx, &session.SessionID{
		ID: sessionID,
	})
	if err != nil {
		return
	}
	userID = int(sess.UserID)
	name = sess.Login

	return
}
