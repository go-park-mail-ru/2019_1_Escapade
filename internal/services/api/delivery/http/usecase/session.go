package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/handler"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
)

// SessionHandler handle requests associated with the session
type SessionHandler struct {
	handler.Handler

	user   api.UserUseCaseI
	auth   infrastructure.AuthService
	trace  infrastructure.ErrorTrace
	logger infrastructure.Logger
}

func NewSessionHandler(
	user api.UserUseCaseI,
	service infrastructure.AuthService,
	trace infrastructure.ErrorTrace,
	logger infrastructure.Logger,
) *SessionHandler {
	return &SessionHandler{
		Handler: *handler.New(logger, trace),

		user:   user,
		auth:   service,
		trace:  trace,
		logger: logger,
	}
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
func (h *SessionHandler) Login(
	rw http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	var (
		user       models.UserPrivateInfo
		publicUser *models.UserPublicInfo
		err        error
		userID     int32
	)

	err = h.GetUser(r, h.auth.HashPassword, &user)
	if err != nil {
		return h.Fail(http.StatusBadRequest, err)
	}

	userID, err = h.user.EnterAccount(
		r.Context(),
		user.Name,
		user.Password,
	)
	if err != nil {
		return h.Fail(
			http.StatusNotFound,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	err = h.auth.CreateToken(rw, user.Name, user.Password)
	if err != nil {
		h.Warning(err, WrnFailedTokenCreate)
	}

	publicUser, err = h.user.FetchOne(r.Context(), userID, 0)
	if err != nil {
		return h.Fail(
			http.StatusInternalServerError,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	return h.Success(http.StatusOK, publicUser)
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
func (h *SessionHandler) Logout(
	rw http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	err := h.auth.DeleteToken(rw, r)
	if err != nil {
		h.Warning(err, "Cant delete token in auth service")
	}

	// TODO send here request to auth service to delete token from database
	return h.Success(http.StatusOK, nil)
}
