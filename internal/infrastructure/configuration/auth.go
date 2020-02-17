package configuration

import (
	"time"
)

// AuthRepository manage getting and setting the configuration of AuthService
type AuthRepository interface {
	Get() Auth
	Set(Auth)
}

// Auth is configuration to initialize infrastructure.AuthService
type Auth struct {
	Salt   string
	Cookie Cookie
	Client AuthClient
}

type AuthClientRepository interface {
	Get() AuthClient
	Set(AuthClient)
}

type AuthClient struct {
	// address of auth service
	Address                string
	ClientID, ClientSecret string
	Scopes                 []string
	RedirectURL            string
}

type CookieRepository interface {
	Get() Cookie
	Set(Cookie)
}

type Cookie struct {
	Path     string
	Lifetime time.Duration
	HTTPOnly bool
	Keys     AuthKeys
}

type AuthKeysRepository interface {
	Get() AuthKeys
	Set(AuthKeys)
}

type AuthKeys struct {
	Access, Type, Refresh, Expire, ReservePrefix string
}
