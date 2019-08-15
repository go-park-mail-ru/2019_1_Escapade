package handlers

import (
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	pg "github.com/vgarvardt/go-oauth2-pg"
	"gopkg.in/oauth2.v3/server"

	"github.com/gorilla/mux"
)

func Router(srv *server.Server, tokenStore *pg.TokenStore) *mux.Router {

	r := mux.NewRouter()

	var auth = r.PathPrefix("/auth").Subrouter()

	auth.Use(mi.Recover)

	auth.HandleFunc("/login", loginHandler)
	auth.HandleFunc("/auth", authHandler)
	auth.HandleFunc("/delete", deleteHandler(srv, tokenStore))
	auth.HandleFunc("/test", testHandler(srv))
	auth.HandleFunc("/token", tokenHandler(srv))
	auth.HandleFunc("/authorize", authorizeHandler(srv))

	return r
}
