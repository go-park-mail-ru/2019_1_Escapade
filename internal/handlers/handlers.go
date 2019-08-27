package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	clients "github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"

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
		rw.Write(b)
	}
}
