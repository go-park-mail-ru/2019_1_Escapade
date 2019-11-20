package handlers

import (
	"gopkg.in/oauth2.v3/server"

	pg "github.com/vgarvardt/go-oauth2-pg"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/gorilla/mux"

	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/router"
)

func Router(srv *server.Server, tokenStore *pg.TokenStore) *mux.Router {

	r := mux.NewRouter()
	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	var auth = r.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/login", loginHandler)
	auth.HandleFunc("/auth", authHandler)
	auth.HandleFunc("/delete", deleteHandler(srv, tokenStore))
	auth.HandleFunc("/test", testHandler(srv))
	auth.HandleFunc("/token", tokenHandler(srv))
	auth.HandleFunc("/authorize", authorizeHandler(srv))

	router.Use(r, mi.Logger)

	return r
}
