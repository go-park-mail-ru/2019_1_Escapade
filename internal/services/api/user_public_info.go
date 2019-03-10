package api

import (
	"errors"
	"escapade/internal/models"
	"net/http"
)

func sendPublicUser(h *Handler, rw http.ResponseWriter, username string, place string) error {

	var (
		user models.UserPublicInfo
		err  error
	)

	if username == "" {
		return errors.New("No username found")
	}

	if user, err = h.DB.GetProfile(username); err != nil {
		return err
	}

	sendSuccessJSON(rw, user, place)
	return err
}
