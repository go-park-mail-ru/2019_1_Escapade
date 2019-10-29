package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"net/http"
)

// Login login
// @Summary login
// @Description login
// @ID Login
// @Success 200 {object} models.UserPublicInfo "Get successfully"
// @Failure 400 {object} models.Result "invalid name or password"
// @Failure 404 {object} models.Result "user not found"
// @Failure 500 {object} models.Result "error with database"
// @Router /session [POST]
func (h *Handler) Login(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "Login"
	var (
		user       models.UserPrivateInfo
		publicUser *models.UserPublicInfo
		err        error
		userID     int32
	)

	err = ih.GetUser(r, h.Auth.Salt, &user)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, err)
	}

	userID, err = h.DB.Login(user.Name, user.Password)
	if err != nil {
		return ih.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	err = auth.CreateTokenInCookies(rw, user.Name, user.Password, h.AuthClient.Config, h.Cookie)
	if err != nil {
		ih.Warning(err, "Cant create token in auth service", place)
	}

	if publicUser, err = h.DB.GetUser(userID, 0); err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, re.NoUserWrapper(err))
	}

	return ih.NewResult(http.StatusOK, place, publicUser, nil)
}

// Logout logout
// @Summary logout
// @Description logout
// @ID Logout
// @Success 200 {object} models.Result "Get successfully"
// @Failure 500 {object} models.Result "server error"
// @Router /session [DELETE]
func (h *Handler) Logout(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "Logout"
	if err := auth.DeleteToken(rw, r, h.Cookie, h.AuthClient); err != nil {
		ih.Warning(err, "Cant delete token in auth service", place)
	}
	return ih.NewResult(http.StatusOK, place, nil, nil)
}

// 85 -> 67
