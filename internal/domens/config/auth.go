package config

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens"
	"golang.org/x/oauth2"
)

type AuthToken struct {
	Auth       Auth
	AuthClient AuthClient
	Cookie     Cookie
}

func NewAuthToken(Auth Auth, AuthClient AuthClient, Cookie Cookie) *AuthToken {
	return &AuthToken{
		Auth:       Auth,
		AuthClient: AuthClient,
		Cookie:     Cookie,
	}
}

// Auth client of auth microservice
//easyjson:json
type Auth struct {
	Salt               string          `json:"salt"`
	AccessTokenExpire  domens.Duration `json:"accessTokenExpire"`
	RefreshTokenExpire domens.Duration `json:"refreshTokenExpire"`
	IsGenerateRefresh  bool            `json:"isGenerateRefresh"`
	WithReserve        bool            `json:"withReserve"`
	TokenType          string          `json:"tokenType"`
	WhiteList          []AuthClient    `json:"whiteList"`
}

//easyjson:json
type AuthClient struct {
	// address of auth service
	Address      string        `json:"address"`
	ClientID     string        `json:"id"`
	ClientSecret string        `json:"secret"`
	Scopes       []string      `json:"scopes"`
	RedirectURL  string        `json:"redirectURL"`
	Config       oauth2.Config `json:"-"`
}

// Cookie set cookie name, path, length, expiration time
// and HTTPonly flag
//easyjson:json
type Cookie struct {
	Path          string     `json:"path"`
	LifetimeHours int        `json:"lifetime_hours"`
	HTTPOnly      bool       `json:"httpOnly"`
	Auth          AuthCookie `json:"keys"`
}

//easyjson:json
type AuthCookie struct {
	AccessToken   string `json:"accessToken"`
	TokenType     string `json:"tokenType"`
	RefreshToken  string `json:"refreshToken"`
	Expire        string `json:"expire"`
	ReservePrefix string `json:"reservePrefix"`
}
