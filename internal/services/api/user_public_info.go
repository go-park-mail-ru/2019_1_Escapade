package api

import (
	"encoding/json"
	"escapade/internal/models"
	"fmt"
	"net/http"
)

func sendPublicUser(h *Handler, rw http.ResponseWriter, username string, place string) (err error) {
	var (
		user  models.UserPublicInfo
		bytes []byte
	)

	if user, err = h.DB.GetProfile(username); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/Me failed")
		return
	}

	if bytes, err = json.Marshal(user); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/Me cant create json")
		return
	}

	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, string(bytes))

	fmt.Println("api/Me ok")
	return
}
