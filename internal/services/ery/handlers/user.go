package eryhandlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	// erydb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"

	"net/http"
	//"github.com/gorilla/mux"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
	//httpSwagger "github.com/swaggo/http-swagger"
)

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
func (h *Handler) CreateUser(rw http.ResponseWriter, r *http.Request) api.Result {
	const place = "CreateUser"

	var user models.User
	err := api.GetUserWithAllFields(r, salt, &user)
	utils.Debug(false, user.Name)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	if err = api.ValidateUser(&user); err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	if err = h.DB.CreateUser(&user); err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, re.UserExistWrapper(err))
	}

	_, err = auth.CreateTokenInHeaders(rw, user.Name, user.Password, h.AuthClient.Config)
	if err != nil {
		api.Warning(err, "Cant create token in auth service", place)
	}

	return api.NewResult(http.StatusCreated, place, &user, nil)
}

// UpdateUser обновить информацию о пользователе(не пароль)
// @Summary update user information
// @Description update public info
// @ID UpdateUser
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid info"
// @Failure 401 {object} models.Result "need authorization"
// @Failure 500 {object} models.Result "error with database"
// @Router /user [PUT]
func (h *Handler) UpdateUser(rw http.ResponseWriter, r *http.Request) api.Result {
	return api.UpdateModel(r, &models.UserUpdate{}, "UpdateProfile", true,
		func(userID int32) (api.JSONtype, error) {
			user, err := h.DB.GetUser(userID)
			return &user, err
		}, func(userI api.JSONtype) error {
			user, ok := userI.(*models.User)
			if !ok {
				return re.NoUpdate()
			}
			return h.DB.UpdateUser(user)
		})
}

func (h *Handler) GetUser(rw http.ResponseWriter, r *http.Request) api.Result {

	const place = "Get user"
	var (
		user   = models.User{}
		err    error
		userID int32
	)

	if userID, err = api.GetUserIDFromAuthRequest(r); err != nil {
		return api.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	if user, err = h.DB.GetUser(userID); err != nil {
		return api.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	return api.NewResult(http.StatusOK, place, &user, nil)
}

func (h *Handler) GetUsers(rw http.ResponseWriter, r *http.Request) api.Result {

	const place = "Get users"
	var (
		users models.Users
		err   error
		name  string
	)

	name = r.FormValue("name")

	if users, err = h.DB.GetUsers(name); err != nil {
		return api.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	return api.NewResult(http.StatusOK, place, &users, nil)
}

func (h *Handler) GetUserByID(rw http.ResponseWriter, r *http.Request, userID int32) api.Result {

	const place = "Get user"
	var (
		user models.User
		err  error
	)

	if user, err = h.DB.GetUser(userID); err != nil {
		return api.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	return api.NewResult(http.StatusOK, place, &user, nil)
}

// 148
