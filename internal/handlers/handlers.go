package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	clients "github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"net/http"
)

// Handler is struct
type Handler struct {
	DB         database.DataBase
	Cookie     config.Cookie
	Clients    *clients.Clients
	AuthClient config.AuthClient
	Auth       config.Auth
}

// JSONtype is interface to be sent by json
type JSONtype interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

type ModelUpdate interface {
	JSONtype

	Update(JSONtype) bool
}

// Result - every handler return it
type Result struct {
	code  int
	place string
	send  JSONtype
	err   error
}

func NewResult(code int, place string, send JSONtype, err error) Result {
	return Result{
		code:  code,
		place: place,
		send:  send,
		err:   err,
	}
}

// TODO добавить группу ожидания и здесь ждать ее, не забыть использовтаь ustils.WaitWithTimeout
func (h *Handler) Close() {
	h.DB.Db.Close()
}

// UpdateModel update any object in DB
// dont use for updating passwords(or hash it before set to BD)
// needAuth - if true, then userID will be taken from auth request
func UpdateModel(r *http.Request, updated ModelUpdate, place string, needAuth bool,
	getFromDB func(userID int32) (JSONtype, error), setToDB func(JSONtype) error) Result {
	var (
		userID int32
		err    error
	)

	// if need userID from auth middleware
	if needAuth {
		if userID, err = GetUserIDFromAuthRequest(r); err != nil {
			return NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
		}
	}

	// updated - new version of object - get it from request
	if err = ModelFromRequest(r, updated); err != nil {
		return NewResult(http.StatusBadRequest, place, nil, err)
	}

	// object - origin object(old version) - get it from bd
	object, err := getFromDB(userID)
	if err != nil {
		return NewResult(http.StatusBadRequest, place, nil, err)
	}

	// update origin object to new version (taking into account empty fields)
	if !updated.Update(object) {
		return NewResult(http.StatusBadRequest, place, nil, re.NoUpdate())
	}

	// try to set updated object to database
	if err = setToDB(object); err != nil {
		return NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return NewResult(http.StatusOK, place, object, nil)
}

func (h *Handler) HandleUser(rw http.ResponseWriter, r *http.Request) {
	var result Result

	switch r.Method {
	case http.MethodPost:
		result = h.CreateUser(rw, r)
	case http.MethodGet:
		result = h.GetMyProfile(rw, r)
	case http.MethodDelete:
		result = h.DeleteUser(rw, r)
	case http.MethodPut:
		result = h.UpdateProfile(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/user wrong request:", r.Method)
	}
	SendResult(rw, result)
	return
}

func (h *Handler) HandleSession(rw http.ResponseWriter, r *http.Request) {
	var result Result

	switch r.Method {
	case http.MethodPost:
		result = h.Login(rw, r)
	case http.MethodDelete:
		result = h.Logout(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	SendResult(rw, result)
	return
}

func (h *Handler) HandleAvatar(rw http.ResponseWriter, r *http.Request) {
	var result Result

	switch r.Method {
	case http.MethodPost:
		result = h.PostImage(rw, r)
	case http.MethodGet:
		result = h.GetImage(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/avatar wrong request:", r.Method)
	}
	SendResult(rw, result)
	return
}

func (h *Handler) HandleOfflineGame(rw http.ResponseWriter, r *http.Request) {
	var result Result

	switch r.Method {
	case http.MethodPost:
		result = h.OfflineSave(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/game wrong request:", r.Method)
	}
	SendResult(rw, result)
	return
}

func (h *Handler) HandleUsersPages(rw http.ResponseWriter, r *http.Request) {
	var result Result

	switch r.Method {
	case http.MethodGet:
		result = h.GetUsers(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/users/pages wrong request:", r.Method)
	}
	SendResult(rw, result)
	return
}

func (h *Handler) HandleUsersPageAmount(rw http.ResponseWriter, r *http.Request) {
	var result Result

	switch r.Method {
	case http.MethodGet:
		result = h.GetUsersPageAmount(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/users/pages_amount wrong request:", r.Method)
	}
	SendResult(rw, result)
	return
}

func SendResult(rw http.ResponseWriter, result Result) {
	if result.code == 0 {
		return
	}

	if result.err != nil {
		sendErrorJSON(rw, result.err, result.place)
	} else {
		sendSuccessJSON(rw, result.send, result.place)
	}
	rw.WriteHeader(result.code)
	Debug(result.err, result.code, result.place)
}

func Debug(catched error, number int, place string) {
	if catched != nil {
		utils.Debug(false, "api/"+place+" failed(code:", number, "). Error message:"+catched.Error())
	} else {
		utils.Debug(false, "api/"+place+" success(code:", number, ")")
	}
}

func Warning(err error, text string, place string) {
	if err != nil {
		utils.Debug(false, "Warning in "+place+".", text, "More:", err.Error())
	} else {
		utils.Debug(false, "Warning in "+place+".", text)
	}
}

// SendErrorJSON send error json
func sendErrorJSON(rw http.ResponseWriter, catched error, place string) {
	result := models.Result{
		Place:   place,
		Success: false,
		Message: catched.Error(),
	}

	if b, err := result.MarshalJSON(); err == nil {
		rw.Write(b)
	}
}

// SendSuccessJSON send object json
func sendSuccessJSON(rw http.ResponseWriter, result JSONtype, place string) {
	if result == nil {
		result = &models.Result{
			Place:   place,
			Success: true,
			Message: "no error",
		}
	}
	if b, err := result.MarshalJSON(); err == nil {
		utils.Debug(false, string(b))
		rw.Write(b)
	}
}
