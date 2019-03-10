package api

import (
	"escapade/internal/models"
	"net/http"
)

func sendPublicUser(h *Handler, rw http.ResponseWriter, username string, place string) (err error) {

	var (
		user models.UserPublicInfo
	)

	defer fixResult(rw, &err, place, user)

	if user, err = h.DB.GetProfile(username); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}
