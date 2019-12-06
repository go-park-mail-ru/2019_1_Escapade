package config

import "golang.org/x/oauth2"

// Auth client of auth microservice
//easyjson:json
type Auth struct {
	Salt               string       `json:"salt"`
	AccessTokenExpire  Duration     `json:"accessTokenExpire"`
	RefreshTokenExpire Duration     `json:"refreshTokenExpire"`
	IsGenerateRefresh  bool         `json:"isGenerateRefresh"`
	WithReserve        bool         `json:"withReserve"`
	TokenType          string       `json:"tokenType"`
	WhiteList          []AuthClient `json:"whiteList"`
}

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

func (conf *Configuration) setOauth2Config() {
	conf.AuthClient.Config = oauth2.Config{
		ClientID:     conf.AuthClient.ClientID,
		ClientSecret: conf.AuthClient.ClientSecret,
		Scopes:       conf.AuthClient.Scopes,
		RedirectURL:  conf.AuthClient.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.AuthClient.Address + "/auth/authorize",
			TokenURL: conf.AuthClient.Address + "/auth/token",
		},
	}
}
