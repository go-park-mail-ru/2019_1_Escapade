package api

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/cookie"

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
		user        models.UserPrivateInfo
		err         error
		found       *models.UserPublicInfo
		sessionName string
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	sessionName = utils.RandomString(16)
	if found, err = h.DB.Login(&user, sessionName); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	token, err := h.Oauth.PasswordCredentialsToken(context.Background(), user.Name, user.Password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		ctx := r.Context()

		sessionID, err := h.Clients.Session.Create(ctx,
			&session.Session{
				UserID: int32(user.ID),
				Login:  user.Name,
			})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("cookie set ", sessionID.ID)
	*/

	//cookie.CreateAndSet(rw, h.Session, sessionName)
	http.SetCookie(rw, cookie.Cookie("access_token", token.AccessToken, token.Expiry))
	http.SetCookie(rw, cookie.Cookie("token_type", token.TokenType, token.Expiry))
	http.SetCookie(rw, cookie.Cookie("refresh_token", token.RefreshToken, token.Expiry))

	utils.SendSuccessJSON(rw, found, place)

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
	const place = "Logout"
	/*
		var (
			err       error
			sessionID string
		)

		if sessionID, err = cookie.GetSessionCookie(r, h.Session); err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
			utils.PrintResult(err, http.StatusUnauthorized, place)
			return
		}
		h.DB.DeleteSession(sessionID)
	*/
	/*
		ctx := context.Background()
		_, err = h.Clients.Session.Delete(ctx,
			&session.SessionID{
				ID: sessionID,
			})
		if err != nil {
			fmt.Println(err)
			return
		}*/

	//cookie.CreateAndSet(rw, h.Session, "")
	http.SetCookie(rw, cookie.Cookie("access_token", "", time.Unix(0, 0)))
	http.SetCookie(rw, cookie.Cookie("token_type", "", time.Unix(0, 0)))
	http.SetCookie(rw, cookie.Cookie("refresh_token", "", time.Unix(0, 0)))

	rw.WriteHeader(http.StatusOK)

	cookie, err := r.Cookie("access_token")
	if err != nil || cookie == nil || cookie.Value == "" {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		utils.Debug(false, "1)something went wrong", err.Error())
		return
	}

	var authServerURL = "http://localhost:9096"
	resp, err := http.Get(fmt.Sprintf("%s/delete?access_token=%s", authServerURL, cookie.Value))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		utils.Debug(false, "2)something went wrong", err.Error())
		return
	}
	defer resp.Body.Close()
	//utils.SendSuccessJSON(rw, nil, place)
	//utils.PrintResult(err, http.StatusOK, place)
	return
}
