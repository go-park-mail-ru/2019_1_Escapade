package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
)

// SessionHandler handle requests associated with the session
type SessionHandler struct {
	ih.Handler
	user database.UserUseCaseI
}

// Init open connections to database
func (h *SessionHandler) Init(c *config.Configuration, db *database.Input) error {
	h.Handler.Init(c)

	h.user = new(database.UserUseCase).Init(db.User, db.Record)
	return h.user.Use(db.Database)
}

// Close connections to database
func (h *SessionHandler) Close() error {
	return h.user.Close()
}

// Handle process any operation associated with user
// authorization: enter and exit
func (h *SessionHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.Login,
		http.MethodDelete:  h.Logout,
		http.MethodOptions: nil})
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

	userID, err = h.user.EnterAccount(user.Name, user.Password)
	if err != nil {
		return ih.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	err = auth.CreateTokenInCookies(rw, user.Name, user.Password, h.AuthClient.Config, h.Cookie)
	if err != nil {
		ih.Warning(err, "Cant create token in auth service", place)
	}

	if publicUser, err = h.user.FetchOne(userID, 0); err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, re.NoUserWrapper(err))
	}

	return ih.NewResult(http.StatusOK, place, publicUser, nil)
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
	const place = "Logout"
	if err := auth.DeleteToken(rw, r, h.Cookie, h.AuthClient); err != nil {
		ih.Warning(err, "Cant delete token in auth service", place)
	}

	// send here request to auth service to delete token from database
	return ih.NewResult(http.StatusOK, place, nil, nil)
}

// 85 -> 67
