package api

import (
	"encoding/json"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"net/http"
)

func getUser(r *http.Request) (user models.UserPrivateInfo, err error) {
	if r.Body == nil {
		err = re.ErrorNoBody()
		return
	}
	_ = json.NewDecoder(r.Body).Decode(&user)
	return
}
