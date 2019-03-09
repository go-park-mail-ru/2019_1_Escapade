package api

import (
	"encoding/json"
	"errors"
	"escapade/internal/models"
	"net/http"
)

func getUser(r *http.Request) (user models.UserPrivateInfo, err error) {
	if r.Body == nil {
		err = errors.New("JSON not found")
		return
	}
	_ = json.NewDecoder(r.Body).Decode(&user)
	/*
		user = models.UserPrivateInfo{

			Name:     r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}*/
	return
}
