package infrastructure

import (
	pg "github.com/vgarvardt/go-oauth2-pg"
	"gopkg.in/oauth2.v3/manage"
)

type TokenManager interface {
	Manager() *manage.Manager
	Store() *pg.TokenStore
	Close() error
}
