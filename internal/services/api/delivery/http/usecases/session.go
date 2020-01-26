package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
)

// SessionHandler handle requests associated with the session
type SessionHandler struct {
	config config.AuthToken
	user   api.UserUseCaseI
}

func NewSessionHandler(c config.AuthToken, user api.UserUseCaseI) *SessionHandler {
	handler := &SessionHandler{
		config: c,
		user:   user,
	}
	return handler
}

// Login login
// @Summary login
// @Description Login into account and get auth2 token.
// @ID Login
// @Tags account
// @Accept  json
// @Param information body models.UserPrivateInfo true "user's name and password"
// @Produce  json
// @Success 200 {object} models.UserPublicInfo "Get successfully"
// @Failure 400 {object} models.Result "invalid name or password"
// @Failure 404 {object} models.Result "Not found"
// @Failure 500 {object} models.Result "Database error"
// @Router /session [POST]
func (h *SessionHandler) Login(rw http.ResponseWriter, r *http.Request) ih.Result {
	var (
		user       models.UserPrivateInfo
		publicUser *models.UserPublicInfo
		err        error
		userID     int32
	)

	err = ih.GetUser(r, h.config.Auth.Salt, &user)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, nil, err)
	}

	userID, err = h.user.EnterAccount(r.Context(), user.Name, user.Password)
	if err != nil {
		return ih.NewResult(http.StatusNotFound, nil, re.NoUserWrapper(err))
	}

	err = auth.CreateTokenInCookies(rw, user.Name, user.Password, h.config.AuthClient.Config, h.config.Cookie)
	if err != nil {
		ih.Warning(err, "Cant create token in auth service")
	}

	publicUser, err = h.user.FetchOne(r.Context(), userID, 0)
	if err != nil {
		return ih.NewResult(http.StatusInternalServerError, nil, re.NoUserWrapper(err))
	}

	return ih.NewResult(http.StatusOK, publicUser, nil)
}

// Logout logout
// @Summary logout
// @Description Logout from account and delete auth2 token.
// @ID Logout
// @Tags account
// @Security OAuth2Application[write]
// @Success 200 {object} models.Result "Get successfully"
// @Failure 500 {object} models.Result "Database error"
// @Router /session [DELETE]
func (h *SessionHandler) Logout(rw http.ResponseWriter, r *http.Request) ih.Result {
	err := auth.DeleteToken(rw, r, h.config.Cookie, h.config.AuthClient)
	if err != nil {
		ih.Warning(err, "Cant delete token in auth service")
	}

	// send here request to auth service to delete token from database
	return ih.NewResult(http.StatusOK, nil, nil)
}
