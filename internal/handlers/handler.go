package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
)

/*
Handler contains all API operations
DB - database, where api work with information
Cookie - cookie settings, more in structure config.Cookie
Clients - grps.Clients, need to connect to Auth server
Auth - auth settinds, more in structure config.Auth
*/
type Handler struct {
	Cookie     config.Cookie
	AuthClient config.AuthClient
	Auth       config.Auth
}

// Init configuration fields
func (h *Handler) Init(c *config.Configuration) {
	h.Cookie = c.Cookie
	h.AuthClient = c.AuthClient
	h.Auth = c.Auth
}
