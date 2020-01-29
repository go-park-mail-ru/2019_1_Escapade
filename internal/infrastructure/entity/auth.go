package entity

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
)

type Auth struct {
	Salt   string
	Cookie config.Cookie
	Auth   config.Auth
	Client config.AuthClient
}
