package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"net/http"
)

// Login login
// @Summary login
// @Description login
// @ID Login
// @Success 200 {object} models.UserPublicInfo "Get successfully"
// @Failure 400 {object} models.Result "invalid name or password"
// @Failure 500 {object} models.Result "server error"
// @Router /session [POST]
func (h *Handler) Login(rw http.ResponseWriter, r *http.Request) {
	const place = "Login"
	var (
		user       models.UserPrivateInfo
		publicUser *models.UserPublicInfo
		err        error
		userID     int32
	)

	user, err = getUser(r)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	userID, err = h.DB.Login(user.Name, user.Password)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	err = auth.CreateToken(rw, h.Oauth, user.Name, user.Password)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
	}

	if publicUser, err = h.DB.GetUser(userID, 0); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}
	utils.SendSuccessJSON(rw, publicUser, place)

	rw.WriteHeader(http.StatusOK)
	utils.PrintResult(err, http.StatusOK, place)
	return
}

// Logout logout
// @Summary logout
// @Description logout
// @ID Logout
// @Success 200 {object} models.Result "Get successfully"
// @Success 401 {object} models.Result "Require authorization"
// @Failure 500 {object} models.Result "server error"
// @Router /session [DELETE]
func (h *Handler) Logout(rw http.ResponseWriter, r *http.Request) {
	if err := auth.DeleteToken(rw, r); err != nil {
		utils.Debug(false, "handler Logout error:", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
	return
}
