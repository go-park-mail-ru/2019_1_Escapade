package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"net/http"

	"context"
	"time"

	ran "math/rand"
)

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
func (h *Handler) GetMyProfile(rw http.ResponseWriter, r *http.Request) Result {

	const place = "GetMyProfile"

	userID, err := GetUserIDFromAuthRequest(r)
	if err != nil {
		return NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	return h.getUser(rw, r, userID, place)
}

// CreateUser godoc
// @Summary create new user
// @Description create new user
// @Accept  json
// @Param name body models.UserPrivateInfo true "User ingo1"
// @Produce  json
// @securitydefinitions.oauth2.password OAuth2Password
// @ID Register
// @Success 201 {object} models.Result "Create user successfully"
// @Header 201 {string} Token "qwerty"
// @Failure 400 {object} models.Result "Invalid information"
// @Router /user [POST]
func (h *Handler) CreateUser(rw http.ResponseWriter, r *http.Request) Result {
	const place = "CreateUser"

	var user models.UserPrivateInfo
	err := GetUserWithAllFields(r, h.Auth.Salt, &user)
	if err != nil {
		return NewResult(http.StatusBadRequest, place, nil, err)
	}

	if err = ValidateUser(&user); err != nil {
		return NewResult(http.StatusBadRequest, place, nil, err)
	}

	if _, err = h.DB.Register(&user); err != nil {
		return NewResult(http.StatusBadRequest, place, nil, re.UserExistWrapper(err))
	}

	err = auth.CreateTokenInCookies(rw, user.Name, user.Password, h.AuthClient.Config, h.Cookie)
	if err != nil {
		Warning(err, "Cant create token in auth service", place)
	}

	return NewResult(http.StatusCreated, place, nil, nil)
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
func (h *Handler) UpdateProfile(rw http.ResponseWriter, r *http.Request) Result {
	const place = "UpdateProfile"

	var (
		user   models.UserPrivateInfo
		err    error
		userID int32
	)

	if err = GetUser(r, h.Auth.Salt, &user); err != nil {
		return NewResult(http.StatusBadRequest, place, nil, err)
	}

	if userID, err = GetUserIDFromAuthRequest(r); err != nil {
		return NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	if err = h.DB.UpdatePlayerPersonalInfo(userID, &user); err != nil {
		return NewResult(http.StatusInternalServerError, place, nil, re.NoUserWrapper(err))
	}

	return NewResult(http.StatusOK, place, nil, nil)
}

// DeleteUser delete account
// @Summary delete account
// @Description delete account
// @ID DeleteAccount
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid input"
// @Failure 500 {object} models.Result "error with database"
// @Router /user [DELETE]
func (h *Handler) DeleteUser(rw http.ResponseWriter, r *http.Request) Result {

	const place = "DeleteUser"
	var (
		user models.UserPrivateInfo
		err  error
	)

	if err = GetUserWithAllFields(r, h.Auth.Salt, &user); err != nil {
		return NewResult(http.StatusBadRequest, place, nil, err)
	}

	if err = h.deleteUserInDB(context.Background(), &user, ""); err != nil {
		return NewResult(http.StatusInternalServerError, place, nil, re.NoUserWrapper(err))
	}

	if err := auth.DeleteToken(rw, r, h.Cookie, h.AuthClient); err != nil {
		Warning(err, "Cant delete token in auth service", place)
	}

	return NewResult(http.StatusOK, place, nil, nil)
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
func (h *Handler) GetProfile(rw http.ResponseWriter, r *http.Request) Result {
	const place = "GetProfile"

	userID, err := h.getUserID(r)
	if err != nil {
		return NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	return h.getUser(rw, r, int32(userID), place)
}

func (h *Handler) getUser(rw http.ResponseWriter, r *http.Request,
	userID int32, place string) Result {

	var (
		err       error
		difficult int
		user      *models.UserPublicInfo
	)

	difficult = h.getDifficult(r)

	if user, err = h.DB.GetUser(userID, difficult); err != nil {
		return NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	photo.GetImages(user)

	return NewResult(http.StatusOK, place, user, nil)
}

func (h *Handler) deleteUserInDB(ctx context.Context,
	user *models.UserPrivateInfo, sessionID string) (err error) {

	if err = h.DB.DeleteAccount(user); err != nil {
		return
	}

	// _, err = h.Clients.Session.Delete(ctx,
	// 	&session.SessionID{
	// 		ID: sessionID,
	// 	})

	return
}

// RandomUsers create n random users
func (h *Handler) RandomUsers(limit int) {

	n := 16
	for i := 0; i < limit; i++ {
		ran.Seed(time.Now().UnixNano())
		user := &models.UserPrivateInfo{
			Name:     utils.RandomString(n),
			Password: utils.RandomString(n)}
		userID, err := h.DB.Register(user)
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
			h.DB.UpdateRecords(int32(userID), record)
		}

	}
}

// 364
