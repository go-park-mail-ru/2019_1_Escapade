package eryhandlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"

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
func (h *Handler) Login(rw http.ResponseWriter, r *http.Request) api.Result {
	const place = "Login"
	var (
		user   *models.User
		err    error
		userID int32
	)

	err = api.GetUser(r, salt, user)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	userID, err = h.DB.GetUserID(user.Name, user.Password)
	if err != nil {
		return api.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	_, err = auth.CreateTokenInHeaders(rw, user.Name, user.Password, h.AuthClient.Config)
	if err != nil {
		api.Warning(err, "Cant create token in auth service", place)
	}

	if user, err = h.DB.GetUser(userID); err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, re.NoUserWrapper(err))
	}

	return api.NewResult(http.StatusOK, place, user, nil)
}

func (h *Handler) UpdatePrivate(rw http.ResponseWriter, r *http.Request) api.Result {
	const place = "UpdatePrivate"
	var (
		user models.UpdatePrivateUser
		err  error
	)

	err = api.ModelFromRequest(r, &user)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}
	user.Old.SetPassword(auth.HashPassword(user.Old.GetPassword(), salt))
	user.New.SetPassword(auth.HashPassword(user.New.GetPassword(), salt))

	if err = h.DB.UpdateUserPrivate(&user); err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, re.NoUserWrapper(err))
	}

	_, err = auth.CreateTokenInHeaders(rw, user.New.Name, user.New.Password, h.AuthClient.Config)
	if err != nil {
		api.Warning(err, "Cant create token in auth service", place)
	}

	return api.NewResult(http.StatusOK, place, nil, nil)
}

// Logout logout
// @Summary logout
// @Description logout
// @ID Logout
// @Success 200 {object} models.Result "Get successfully"
// @Failure 500 {object} models.Result "server error"
// @Router /session [DELETE]
func (h *Handler) Logout(rw http.ResponseWriter, r *http.Request) api.Result {
	const place = "Logout"
	if err := auth.DeleteFromHeader(rw, h.AuthClient); err != nil {
		api.Warning(err, "Cant delete token in auth service", place)
	}
	return api.NewResult(http.StatusOK, place, nil, nil)
}

// 85 -> 67
