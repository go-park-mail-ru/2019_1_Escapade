package api

import (
	"escapade/internal/models"
	"fmt"
	"net/http"
)

func sendPublicUser(h *Handler, rw http.ResponseWriter, username string, place string) error {

	var (
		user models.UserPublicInfo
		err  error
	)

	if user, err = h.DB.GetProfile(username); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println(place + " failed")
		return err
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, user, place)
	fmt.Println(place + " ok")
	return err
}
