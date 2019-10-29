package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"

	clients "github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"

	"net/http"
)

/*
Handler contains all API operations
DB - database, where api work with information
Cookie - cookie settings, more in structure config.Cookie
Clients - grps.Clients, need to connect to Auth server
Auth - auth settinds, more in structure config.Auth
*/
type Handler struct {
	DB         database.DataBase
	Cookie     config.Cookie
	Clients    *clients.Clients
	AuthClient config.AuthClient
	Auth       config.Auth
}

// TODO добавить группу ожидания и здесь ждать ее, не забыть использовтаь ustils.WaitWithTimeout
func (h *Handler) Close() {
	h.DB.Db.Close()
}

// HandleUser process any operation associated with user
// profile: create, receive, update, and delete
func (h *Handler) HandleUser(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.CreateUser,
		http.MethodGet:     h.GetMyProfile,
		http.MethodDelete:  h.DeleteUser,
		http.MethodPut:     h.UpdateProfile,
		http.MethodOptions: nil})
}

// HandleSession process any operation associated with user
// authorization: enter and exit
func (h *Handler) HandleSession(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.Login,
		http.MethodDelete:  h.Logout,
		http.MethodOptions: nil})
}

// TODO add deleting
// HandleAvatar process any operation associated with user
// avatar: load and get
func (h *Handler) HandleAvatar(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.PostImage,
		http.MethodGet:     h.GetImage,
		http.MethodOptions: nil})
}

// TODO add receipt
// HandleOfflineGame process any operation associated with offline
// games: save
func (h *Handler) HandleOfflineGame(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.OfflineSave,
		http.MethodOptions: nil})
}

// HandleUsersPages process any operation associated with users
// list: receive
func (h *Handler) HandleUsersPages(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet:     h.GetUsers,
		http.MethodOptions: nil})
}

// HandleUsersPageAmount process any operation associated with
// amount of pages in user list: receive
func (h *Handler) HandleUsersPageAmount(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet:     h.GetUsersPageAmount,
		http.MethodOptions: nil})
}

// 128 -> 88
