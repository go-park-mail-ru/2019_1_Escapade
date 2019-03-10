package api

import (
	"escapade/internal/models"
	"net/http"
)

func sendPublicUser(h *Handler, rw http.ResponseWriter, username string, place string) error {

	var (
		user models.UserPublicInfo
		err  *error
	)

	defer fixResult(rw, err, place, user)

	if user, *err = h.DB.GetProfile(username); *err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return *err
	}

	rw.WriteHeader(http.StatusOK)
	return *err
}
