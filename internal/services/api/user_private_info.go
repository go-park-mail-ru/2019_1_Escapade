package api

import (
	"escapade/internal/models"
	"net/http"
)

func getUser(r *http.Request) models.UserPrivateInfo {
	user := models.UserPrivateInfo{
		Name:     r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	return user
}
