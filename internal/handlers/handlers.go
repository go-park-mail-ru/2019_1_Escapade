package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	clients "github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"

	"net/http"

	"golang.org/x/oauth2"
)

// Handler is struct
type Handler struct {
	DB        database.DataBase
	Session   config.SessionConfig
	WebSocket config.WebSocketSettings
	Game      config.GameConfig
	Clients   *clients.Clients
	Oauth     oauth2.Config
}

// TODO добавить группу ожидания и здесь ждать ее, не забыть использовтаь ustils.WaitWithTimeout
func (h *Handler) Close() {
	h.DB.Db.Close()
}

func (h *Handler) HandleUser(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		h.CreateUser(rw, r)
	case http.MethodGet:
		h.GetMyProfile(rw, r)
	case http.MethodDelete:
		h.DeleteUser(rw, r)
	case http.MethodPut:
		h.UpdateProfile(rw, r)
	case http.MethodOptions:
		rw.WriteHeader(http.StatusOK)
		utils.Debug(false, "cors:", rw.Header().Get("Access-Control-Allow-Methods"))

		return
	default:
		utils.Debug(true, "/api/user wrong request:", r.Method)
	}
	return
}

func (h *Handler) HandleSession(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Login(rw, r)
	case http.MethodDelete:
		h.Logout(rw, r)
	case http.MethodOptions:
		utils.Debug(false, "cors Access-Control-Allow-Methods:", rw.Header().Get("Access-Control-Allow-Methods"))
		return
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	return
}

func (h *Handler) HandleAvatar(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.PostImage(rw, r)
	case http.MethodGet:
		h.GetImage(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/avatar wrong request:", r.Method)
	}
	return
}

func (h *Handler) HandleOfflineGame(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.OfflineSave(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/game wrong request:", r.Method)
	}
	return
}

func (h *Handler) HandleUsersPages(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUsers(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/users/pages wrong request:", r.Method)
	}
	return
}
