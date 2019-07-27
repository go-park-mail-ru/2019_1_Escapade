package api

import (
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/cookie"

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
func (h *Handler) GetMyProfile(rw http.ResponseWriter, r *http.Request) {

	const place = "GetMyProfile"
	var (
		err    error
		userID int
	)

	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	h.getUser(rw, r, userID)

	return
}

// CreateUser godoc
// @Summary create new user
// @Description create new user
// @ID Register
// @Success 201 {object} models.Result "Create user successfully"
// @Header 201 {string} Token "qwerty"
// @Failure 400 {object} models.Result "Invalid information"
// @Router /user [POST]
func (h *Handler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateUser"
	var (
		user      models.UserPrivateInfo
		err       error
		sessionID string
	)

	if user, err = getUserWithAllFields(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if err = validateUser(&user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
	}

	if _, sessionID, err = h.createUserInDB(r.Context(), user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	cookie.CreateAndSet(rw, h.Session, sessionID)
	rw.WriteHeader(http.StatusCreated)
	utils.SendSuccessJSON(rw, nil, place)
	utils.PrintResult(err, http.StatusCreated, place)
	return
}

// UpdateProfile godoc
// @Summary update user information
// @Description update public info
// @ID UpdateProfile
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid info"
// @Failure 401 {object} models.Result "need authorization"
// @Router /user [PUT]
func (h *Handler) UpdateProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateProfile"

	var (
		user   models.UserPrivateInfo
		err    error
		userID int
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	if err = h.DB.UpdatePlayerPersonalInfo(userID, &user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	utils.SendSuccessJSON(rw, nil, place)
	utils.PrintResult(err, http.StatusOK, place)
	return
}

// DeleteUser delete account
// @Summary delete account
// @Description delete account
// @ID DeleteAccount
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid input"
// @Failure 500 {object} models.Result "server error"
// @Router /user [DELETE]
func (h *Handler) DeleteUser(rw http.ResponseWriter, r *http.Request) {

	const place = "DeleteUser"
	var (
		user models.UserPrivateInfo
		err  error
	)

	if user, err = getUserWithAllFields(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if err = h.deleteUserInDB(context.Background(), &user, ""); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	cookie.CreateAndSet(rw, h.Session, "")
	rw.WriteHeader(http.StatusOK)
	utils.SendSuccessJSON(rw, nil, place)
	utils.PrintResult(err, http.StatusOK, place)
	return
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
func (h *Handler) GetProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "GetProfile"

	var (
		err    error
		userID int
	)

	if userID, err = h.getUserID(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	h.getUser(rw, r, userID)
	return
}

func (h *Handler) getUser(rw http.ResponseWriter, r *http.Request, userID int) {
	const place = "GetProfile"

	var (
		err       error
		difficult int
		user      *models.UserPublicInfo
	)

	difficult = h.getDifficult(r)

	if user, err = h.DB.GetUser(userID, difficult); err != nil {

		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}
	photo.GetImages(user)

	utils.SendSuccessJSON(rw, user, place)

	rw.WriteHeader(http.StatusOK)
	utils.PrintResult(err, http.StatusOK, place)
	return
}

func (h *Handler) createUserInDB(ctx context.Context,
	user models.UserPrivateInfo) (userID int, sessionID string, err error) {

	sessionID = utils.RandomString(16)
	if userID, err = h.DB.Register(&user, sessionID); err != nil {
		return
	}

	/*sessID, err := h.Clients.Session.Create(ctx,
		&session.Session{
			UserID: int32(userID),
			Login:  user.Name,
		})

	if err != nil {
		return
	}

	sessionID = sessID.ID*/
	return
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
		user := models.UserPrivateInfo{
			Name:     utils.RandomString(n),
			Password: utils.RandomString(n)}
		userID, _, err := h.createUserInDB(context.Background(), user)
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
			h.DB.UpdateRecords(userID, record)
		}

	}
}

func validateUser(user *models.UserPrivateInfo) error {
	name := strings.TrimSpace(user.Name)
	if name == "" || len(name) < 3 {
		return re.ErrorInvalidName()
	}
	user.Name = name

	password := strings.TrimSpace(user.Password)
	if len(password) < 3 {
		return re.ErrorInvalidPassword()
	}
	return nil
}
