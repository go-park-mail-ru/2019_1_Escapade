package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func sendPublicUser(h *Handler, rw http.ResponseWriter, username string, place string) (err error) {
	user, erro := h.DB.GetProfile(username)
	err = erro
	if err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/Me failed")
		return
	}

	bytes, errJSON := json.Marshal(user)
	if errJSON == nil {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintln(rw, string(bytes))

		fmt.Println("api/Me ok")
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/Me cant create json")
	}
	return
}
