package handlers

import (
	"strconv"
	"net/http"
	"context"
	"time"
	ran "math/rand"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/auth"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
)

type UserHandler struct {
	ih.Handler
	user   database.UserUseCaseI
	record database.RecordUseCaseI
}

func (h *UserHandler) Init(c *config.Configuration, DB idb.DatabaseI,
	userDB database.UserRepositoryI, recordDB database.RecordRepositoryI) error {
	h.Handler.Init(c)

	h.user = &database.UserUseCase{}
	h.user.Init(userDB, recordDB)
	err := h.user.Use(DB)
	if err != nil {
		return err
	}

	h.record = &database.RecordUseCase{}
	h.record.Init(recordDB)
	err = h.record.Use(DB)
	if err != nil {
		return err
	}
	return nil
}

func (h *UserHandler) Close() {
	h.user.Close()
	h.record.Close()
}

// Handle process any operation associated with user
// profile: create, receive, update, and delete
func (h *UserHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.CreateUser,
		http.MethodGet:     h.GetMyProfile,
		http.MethodDelete:  h.DeleteUser,
		http.MethodPut:     h.UpdateProfile,
		http.MethodOptions: nil})
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
func (h *UserHandler) GetMyProfile(rw http.ResponseWriter, r *http.Request) ih.Result {

	const place = "GetMyProfile"

	userID, err := ih.GetUserIDFromAuthRequest(r)
	if err != nil {
		return ih.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	return h.getUser(rw, r, userID, place)
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
func (h *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "CreateUser"

	var user models.UserPrivateInfo
	err := ih.GetUserWithAllFields(r, h.Auth.Salt, &user)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, err)
	}

	if err = ih.ValidateUser(&user); err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, err)
	}

	if _, err = h.user.CreateAccount(&user); err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.UserExistWrapper(err))
	}

	err = auth.CreateTokenInCookies(rw, user.Name, user.Password, h.AuthClient.Config, h.Cookie)
	if err != nil {
		ih.Warning(err, "Cant create token in auth service", place)
	}

	return ih.NewResult(http.StatusCreated, place, nil, nil)
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
func (h *UserHandler) UpdateProfile(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "UpdateProfile"

	var (
		user   models.UserPrivateInfo
		err    error
		userID int32
	)

	if err = ih.GetUser(r, h.Auth.Salt, &user); err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, err)
	}

	if userID, err = ih.GetUserIDFromAuthRequest(r); err != nil {
		return ih.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	if err = h.user.UpdateAccount(userID, &user); err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, re.NoUserWrapper(err))
	}

	return ih.NewResult(http.StatusOK, place, nil, nil)
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
func (h *UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) ih.Result {

	const place = "DeleteUser"
	var (
		user models.UserPrivateInfo
		err  error
	)

	if err = ih.GetUserWithAllFields(r, h.Auth.Salt, &user); err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, err)
	}

	if err = h.deleteUserInDB(context.Background(), &user, ""); err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, re.NoUserWrapper(err))
	}

	if err := auth.DeleteToken(rw, r, h.Cookie, h.AuthClient); err != nil {
		ih.Warning(err, "Cant delete token in auth service", place)
	}

	return ih.NewResult(http.StatusOK, place, nil, nil)
}

func (h *UserHandler) getUser(rw http.ResponseWriter, r *http.Request,
	userID int32, place string) ih.Result {

	var (
		err       error
		difficult string
		user      *models.UserPublicInfo
	)

	difficult = getDifficult(r)
	difficultI, err := strconv.Atoi(difficult)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.NoUserWrapper(err))
	}

	if user, err = h.user.FetchOne(userID, difficultI); err != nil {
		return ih.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	photo.GetImages(user)

	return ih.NewResult(http.StatusOK, place, user, nil)
}

func (h *UserHandler) deleteUserInDB(ctx context.Context,
	user *models.UserPrivateInfo, sessionID string) (err error) {

	if err = h.user.DeleteAccount(user); err != nil {
		return
	}

	// _, err = h.Clients.Session.Delete(ctx,
	// 	&session.SessionID{
	// 		ID: sessionID,
	// 	})

	return
}

// RandomUsers create n random users
func (h *UserHandler) RandomUsers(limit int) {

	n := 16
	for i := 0; i < limit; i++ {
		ran.Seed(time.Now().UnixNano())
		user := &models.UserPrivateInfo{
			Name:     utils.RandomString(n),
			Password: utils.RandomString(n)}
		userID, err := h.user.CreateAccount(user)
		if err != nil {
			utils.Debug(true, "cant register random")
			return
		}

		for j := 0; j < 4; j++ {
			record := &models.Record{
				Score:       ran.Intn(1000000),
				Time:        float64(ran.Intn(10000)),
				Difficult:   j,
				SingleTotal: ran.Intn(2),
				OnlineTotal: ran.Intn(2),
				SingleWin:   ran.Intn(2),
				OnlineWin:   ran.Intn(2)}
			h.record.Update(int32(userID), record)
		}

	}
}

// 364
