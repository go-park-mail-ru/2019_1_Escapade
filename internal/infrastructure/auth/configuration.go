package auth

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens"
	"time"
)

type Configuration struct {
	Cookie Cookie
	Auth   Auth
	Client AuthClient
}

type ConfigurationJSON struct {
	Cookie CookieJSON     `json:"session"`
	Auth   AuthJSON       `json:"auth"`
	Client AuthClientJSON `json:"authClient"`
}

func (c ConfigurationJSON) Get() Configuration {
	return Configuration{
		Cookie: c.Cookie.Get(),
		Auth:   c.Auth.Get(),
		Client: c.Client.Get(),
	}
}

type Auth struct {
	Salt               string
	AccessTokenExpire  time.Duration
	RefreshTokenExpire time.Duration
	IsGenerateRefresh  bool
	WithReserve        bool
	TokenType          string
	WhiteList          []AuthClient
}

// Auth client of auth microservice
//easyjson:json
type AuthJSON struct {
	Salt               string           `json:"salt"`
	AccessTokenExpire  domens.Duration  `json:"accessTokenExpire"`
	RefreshTokenExpire domens.Duration  `json:"refreshTokenExpire"`
	IsGenerateRefresh  bool             `json:"isGenerateRefresh"`
	WithReserve        bool             `json:"withReserve"`
	TokenType          string           `json:"tokenType"`
	WhiteList          []AuthClientJSON `json:"whiteList"` // ! never used
}

func (c AuthJSON) Get() Auth {
	var newWhiteList = make([]AuthClient, len(c.WhiteList))
	for i, e := range c.WhiteList {
		newWhiteList[i] = e.Get()
	}
	return Auth{
		Salt:               c.Salt,
		AccessTokenExpire:  c.AccessTokenExpire.Duration,
		RefreshTokenExpire: c.RefreshTokenExpire.Duration,
		IsGenerateRefresh:  c.IsGenerateRefresh,
		WithReserve:        c.WithReserve,
		TokenType:          c.TokenType,
		WhiteList:          newWhiteList,
	}
}

type AuthClient struct {
	// address of auth service
	Address      string
	ClientID     string
	ClientSecret string
	Scopes       []string
	RedirectURL  string
}

//easyjson:json
type AuthClientJSON struct {
	// address of auth service
	Address      string   `json:"address"`
	ClientID     string   `json:"id"`
	ClientSecret string   `json:"secret"`
	Scopes       []string `json:"scopes"`
	RedirectURL  string   `json:"redirectURL"`
}

func (c AuthClientJSON) Get() AuthClient {
	return AuthClient(c)
}

type Cookie struct {
	Path          string
	LifetimeHours int
	HTTPOnly      bool
	Auth          AuthCookie
}

// Cookie set cookie name, path, length, expiration time
// and HTTPonly flag
//easyjson:json
type CookieJSON struct {
	Path          string         `json:"path"`
	LifetimeHours int            `json:"lifetime_hours"`
	HTTPOnly      bool           `json:"httpOnly"`
	Auth          AuthCookieJSON `json:"keys"`
}

func (c CookieJSON) Get() Cookie {
	return Cookie{
		Path:          c.Path,
		LifetimeHours: c.LifetimeHours,
		HTTPOnly:      c.HTTPOnly,
		Auth:          c.Auth.Get(),
	}
}

type AuthCookie struct {
	AccessToken   string
	TokenType     string
	RefreshToken  string
	Expire        string
	ReservePrefix string
}

//easyjson:json
type AuthCookieJSON struct {
	AccessToken   string `json:"accessToken"`
	TokenType     string `json:"tokenType"`
	RefreshToken  string `json:"refreshToken"`
	Expire        string `json:"expire"`
	ReservePrefix string `json:"reservePrefix"`
}

func (c AuthCookieJSON) Get() AuthCookie {
	return AuthCookie(c)
}
