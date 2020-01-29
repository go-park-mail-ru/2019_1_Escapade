package handlers

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	delivery "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http"
)

// UserHandler handle requests associated with the user
type UserHandler struct {
	auth  infrastructure.AuthService
	photo infrastructure.PhotoServiceI

	user   api.UserUseCaseI
	record api.RecordUseCaseI
	rep    delivery.RepositoryI

	trace infrastructure.ErrorTrace
	log   infrastructure.LoggerI
}

func NewUserHandler(
	user api.UserUseCaseI,
	record api.RecordUseCaseI,
	rep delivery.RepositoryI,
	auth infrastructure.AuthService,
	photo infrastructure.PhotoServiceI,
	trace infrastructure.ErrorTrace,
	log infrastructure.LoggerI,
) *UserHandler {
	handler := &UserHandler{
		auth:  auth,
		photo: photo,

		user:   user,
		record: record,
		rep:    rep,

		trace: trace,
		log:   log,
	}
	return handler
}

// GetMyProfile get user public information
// @Summary get user public information
// @Description  get user's best score and best time for a given difficulty, user's id, name and photo of current user. The current one is the one whose token is provided.
// @ID GetMyProfile
// @Security OAuth2Application[read]
// @Tags account
// @Accept  json
// @Param difficult query int false "which difficult records will be given" default(0)
// @Produce  json
// @Success 200 {object} models.UserPublicInfo "Get user successfully"
// @Failure 401 {object} models.Result "Authorization required"
// @Router /user [GET]
func (h *UserHandler) GetMyProfile(
	rw http.ResponseWriter,
	r *http.Request,
) ih.Result {

	userID, err := ih.GetUserIDFromAuthRequest(r)
	if err != nil {
		return ih.NewResult(
			http.StatusUnauthorized,
			nil,
			h.trace.WrapWithText(err, ErrAuth),
		)
	}

	return h.getUser(rw, r, userID)
}

// CreateUser create user
// @Summary create new user
// @Description create new account and get oauth2 token
// @ID CreateUser
// @Tags account
// @Accept  json
// @Param information body models.UserPrivateInfo true "user's name and password"
// @Produce  json
// @Success 201 {object} models.Result "Create user successfully"
// @Failure 400 {object} models.Result "Invalid information"
// @Router /user [POST]
func (h *UserHandler) CreateUser(
	rw http.ResponseWriter,
	r *http.Request,
) ih.Result {
	var user models.UserPrivateInfo
	err := ih.GetUserWithAllFields(r, h.auth, &user)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, nil, err)
	}

	if err = ih.ValidateUser(&user); err != nil {
		return ih.NewResult(http.StatusBadRequest, nil, err)
	}

	_, err = h.user.CreateAccount(r.Context(), &user)
	if err != nil {
		return ih.NewResult(
			http.StatusBadRequest,
			nil,
			h.trace.WrapWithText(err, ErrUserAlreadyExist),
		)
	}

	// TODO а может тут вызывать какой нибудь метод обработчика сессий?
	err = h.auth.CreateToken(rw, user.Name, user.Password)
	if err != nil {
		ih.Warning(err, "Cant create token in auth service")
	}

	return ih.NewResult(http.StatusCreated, nil, nil)
}

// UpdateProfile update user's name or password
// @Summary update user sensitive data
// @Description update name or/and password of current user. The current one is the one whose token is provided.
// @ID UpdateProfile
// @Security OAuth2Application[write]
// @Tags account
// @Accept  json
// @Param information body models.UserPrivateInfo true "user's name and password"
// @Produce  json
// @Success 200 {object} models.Result "Update successfully"
// @Failure 400 {object} models.Result "Invalid data for update"
// @Failure 401 {object} models.Result "Authorization required"
// @Failure 500 {object} models.Result "Database error"
// @Router /user [PUT]
func (h *UserHandler) UpdateProfile(
	rw http.ResponseWriter,
	r *http.Request,
) ih.Result {
	var (
		user   models.UserPrivateInfo
		err    error
		userID int32
	)

	if err = ih.GetUser(r, h.auth, &user); err != nil {
		return ih.NewResult(http.StatusBadRequest, nil, err)
	}

	userID, err = ih.GetUserIDFromAuthRequest(r)
	if err != nil {
		return ih.NewResult(
			http.StatusUnauthorized,
			nil,
			h.trace.WrapWithText(err, ErrAuth),
		)
	}

	err = h.user.UpdateAccount(r.Context(), userID, &user)
	if err != nil {
		return ih.NewResult(
			http.StatusInternalServerError,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	return ih.NewResult(http.StatusOK, nil, nil)
}

// DeleteUser delete account
// @Summary delete account
// @Description delete account
// @ID DeleteAccount
// @Tags account
// @Accept  json
// @Param information body models.UserPrivateInfo true "user's name and password.  You are required to pass in the body of the request user name and password to confirm that you are the owner of the account."
// @Produce  json
// @Success 200 {object} models.Result "Delete successfully"
// @Failure 400 {object} models.Result "Invalid data for delete"
// @Failure 500 {object} models.Result "Database error"
// @Router /user [DELETE]
func (h *UserHandler) DeleteUser(
	rw http.ResponseWriter,
	r *http.Request,
) ih.Result {
	var (
		user models.UserPrivateInfo
		err  error
	)

	err = ih.GetUserWithAllFields(r, h.auth, &user)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, nil, err)
	}

	err = h.deleteUserInDB(context.Background(), &user, "")
	if err != nil {
		return ih.NewResult(
			http.StatusInternalServerError,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	err = h.auth.DeleteToken(rw, r)
	if err != nil {
		ih.Warning(err, WrnFailedTokenDelete)
	}

	return ih.NewResult(http.StatusOK, nil, nil)
}

func (h *UserHandler) getUser(
	rw http.ResponseWriter,
	r *http.Request,
	userID int32,
) ih.Result {

	var (
		err       error
		difficult string
		user      *models.UserPublicInfo
	)

	difficult = h.rep.GetDifficult(r)
	difficultI, err := strconv.Atoi(difficult)
	if err != nil {
		return ih.NewResult(
			http.StatusBadRequest,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	user, err = h.user.FetchOne(r.Context(), userID, difficultI)
	if err != nil {
		return ih.NewResult(
			http.StatusNotFound,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	h.photo.GetImages(user)

	return ih.NewResult(http.StatusOK, user, nil)
}

func (h *UserHandler) deleteUserInDB(
	ctx context.Context,
	user *models.UserPrivateInfo,
	sessionID string,
) error {

	err := h.user.DeleteAccount(ctx, user)
	if err != nil {
		return err
	}

	// _, err = h.Clients.Session.Delete(ctx,
	// 	&session.SessionID{
	// 		ID: sessionID,
	// 	})

	return nil
}

// RandomUsers create n random users
func (h *UserHandler) RandomUsers(
	c context.Context,
	limit int,
	timeout time.Duration,
) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()
	n := 16
	for i := 0; i < limit; i++ {
		rand.Seed(time.Now().UnixNano())
		user := &models.UserPrivateInfo{
			Name:     utils.RandomString(n),
			Password: utils.RandomString(n),
		}
		userID, err := h.user.CreateAccount(ctx, user)
		if err != nil {
			h.log.Fatalf("cant register random")
			return
		}

		for j := 0; j < 4; j++ {
			h.record.Update(ctx, int32(userID), &models.Record{
				Score:       rand.Intn(1000000),
				Time:        float64(rand.Intn(10000)),
				Difficult:   j,
				SingleTotal: rand.Intn(2),
				OnlineTotal: rand.Intn(2),
				SingleWin:   rand.Intn(2),
				OnlineWin:   rand.Intn(2)})
		}
	}
}
