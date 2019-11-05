package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"net/http"

	"context"
	"time"

	ran "math/rand"
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

// GetMyProfile godoc
// @Summary get public information about that user
// @Description get public information about that user
// @ID GetMyProfile
// @Produce  json
// @Param enumstring query string false "string enums" Enums(A, B, C)
// @Param enumint query int false "int enums" Enums(1, 2, 3)
// @Param enumnumber query number false "int enums" Enums(1.1, 1.2, 1.3)
// @Param string query string false "string valid" minlength(5) maxlength(10)
// @Param int query int false "int valid" mininum(1) maxinum(10)
// @Param default query string false "string default" default(A)
// @Success 201 {object} models.Result "Create user successfully"
// @Header 201 {string} Token "qwerty"
// @Failure 401 {object} models.Result "Invalid information"
// @Router /user [GET]
func (h *UserHandler) GetMyProfile(rw http.ResponseWriter, r *http.Request) ih.Result {

	const place = "GetMyProfile"

	userID, err := ih.GetUserIDFromAuthRequest(r)
	if err != nil {
		return ih.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	return h.getUser(rw, r, userID, place)
}

// CreateUser godoc
// @Summary create new user
// @Description create new user
// @Accept  json
// @Param name body models.UserPrivateInfo true "User info1"
// @Produce  json
// @securitydefinitions.oauth2.password OAuth2Password
// @ID Register
// @Success 201 {object} models.Result "Create user successfully"
// @Header 201 {string} Token "qwerty"
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

// UpdateProfile godoc
// @Summary update user information
// @Description update public info
// @ID UpdateProfile
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid info"
// @Failure 401 {object} models.Result "need authorization"
// @Failure 500 {object} models.Result "error with database"
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
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid input"
// @Failure 500 {object} models.Result "error with database"
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

// GetProfile godoc
// @Summary Get public user inforamtion
// @Description get user's best score and best time for a given difficulty, user's id, name and photo
// @ID GetProfile
// @Accept  json
// @Produce  json
// @Param name path string false "User name"
// @Success 200 {object} models.UserPublicInfo "Profile found successfully"
// @Failure 400 {object} models.Result "Invalid username"
// @Failure 404 {object} models.Result "User not found"
// @Router /users/{name}/profile [GET]
func (h *UserHandler) GetProfile(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "GetProfile"

	userID, err := getUserID(r)
	if err != nil {
		return ih.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	return h.getUser(rw, r, int32(userID), place)
}

func (h *UserHandler) getUser(rw http.ResponseWriter, r *http.Request,
	userID int32, place string) ih.Result {

	var (
		err       error
		difficult int
		user      *models.UserPublicInfo
	)

	difficult = getDifficult(r)

	if user, err = h.user.FetchOne(userID, difficult); err != nil {
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
