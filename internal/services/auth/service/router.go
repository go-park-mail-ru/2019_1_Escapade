package service

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/handlers"
	pg "github.com/vgarvardt/go-oauth2-pg"
	"gopkg.in/oauth2.v3/server"
)

type Router struct {
	srv        *server.Server
	tokenStore *pg.TokenStore
}

func (r *Router) Init(srv *server.Server, store *pg.TokenStore) *Router {
	r.srv = srv
	r.tokenStore = store
	return r
}

func (r *Router) Router() http.Handler {
	return handlers.Router(r.srv, r.tokenStore)
}
