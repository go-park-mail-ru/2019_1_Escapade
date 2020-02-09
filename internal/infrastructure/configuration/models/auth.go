package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

//easyjson:json
type Auth struct {
	Salt            string          `json:"salt"`
	Cookie          Cookie          `json:"cookie"`
	TokenGeneration TokenGeneration `json:"token_generation"`
	Client          AuthClient      `json:"auth_client"`
}

//easyjson:json
type TokenGeneration struct {
	AccessExpire      models.Duration `json:"accessTokenExpire"`
	RefreshExpire     models.Duration `json:"refreshTokenExpire"`
	IsGenerateRefresh bool            `json:"isGenerateRefresh"`
	TokenType         string          `json:"tokenType"`
}

//easyjson:json
type AuthClient struct {
	// address of auth service
	Address      string   `json:"address" env:"auth_address"`
	ClientID     string   `json:"id"`
	ClientSecret string   `json:"secret"`
	Scopes       []string `json:"scopes"`
	RedirectURL  string   `json:"redirect_url" env:"auth_redirect_url"`
}

// Cookie set cookie name, path, length, expiration time
// and HTTPonly flag
//easyjson:json
type Cookie struct {
	Path     string          `json:"path"`
	Lifetime models.Duration `json:"lifetime"`
	HTTPOnly bool            `json:"http_only"`
	Keys     AuthKeys        `json:"keys"`
}

type AuthKeys struct {
	Access  string `json:"access"`
	Type    string `json:"type"`
	Refresh string `json:"refresh"`
	Expire  string `json:"expire"`
}

// Get configuration.Auth from json model
// implementation of AuthRepository
func (a *Auth) Get() configuration.Auth {
	return configuration.Auth{
		Salt:            a.Salt,
		Cookie:          a.Cookie.Get(),
		TokenGeneration: a.TokenGeneration.Get(),
		Client:          a.Client.Get(),
	}
}

// Set data from configuration.Auth
// implementation of AuthRepository
func (a *Auth) Set(c configuration.Auth) {
	a.Salt = c.Salt
	a.Cookie.Set(c.Cookie)
	a.TokenGeneration.Set(c.TokenGeneration)
	a.Client.Set(c.Client)
}

// Get configuration.TokenGeneration from json model
// implementation of TokenGenerationRepository
func (a *TokenGeneration) Get() configuration.TokenGeneration {
	return configuration.TokenGeneration{
		AccessExpire:      a.AccessExpire.Duration,
		RefreshExpire:     a.RefreshExpire.Duration,
		IsGenerateRefresh: a.IsGenerateRefresh,
		TokenType:         a.TokenType,
	}
}

// Set data from configuration.TokenGeneration
// implementation of TokenGenerationRepository
func (a *TokenGeneration) Set(c configuration.TokenGeneration) {
	a.AccessExpire.Init(c.AccessExpire)
	a.RefreshExpire.Init(c.RefreshExpire)
	a.IsGenerateRefresh = c.IsGenerateRefresh
	a.TokenType = c.TokenType
}

// Get configuration.AuthClient from json model
// implementation of AuthClientRepository
func (a *AuthClient) Get() configuration.AuthClient {
	return configuration.AuthClient{
		Address:      a.Address,
		ClientID:     a.ClientID,
		ClientSecret: a.ClientSecret,
		Scopes:       a.Scopes,
		RedirectURL:  a.RedirectURL,
	}
}

// Set data from configuration.AuthClient
// implementation of AuthClientRepository
func (a *AuthClient) Set(c configuration.AuthClient) {
	a.Address = c.Address
	a.ClientID = c.ClientID
	a.ClientSecret = c.ClientSecret
	a.Scopes = c.Scopes
	a.RedirectURL = c.RedirectURL
}

// Get configuration.Cookie from json model
// implementation of CookieRepository
func (a *Cookie) Get() configuration.Cookie {
	return configuration.Cookie{
		Path:     a.Path,
		Lifetime: a.Lifetime.Duration,
		HTTPOnly: a.HTTPOnly,
		Keys:     a.Keys.Get(),
	}
}

// Set data from configuration.Cookie
// implementation of CookieRepository
func (a *Cookie) Set(c configuration.Cookie) {
	a.Path = c.Path
	a.Lifetime.Init(c.Lifetime)
	a.HTTPOnly = c.HTTPOnly
	a.Keys.Set(c.Keys)
}

// Get configuration.AuthKeys from json model
// implementation of AuthKeysRepository
func (a *AuthKeys) Get() configuration.AuthKeys {
	return configuration.AuthKeys{
		Access:  a.Access,
		Type:    a.Type,
		Refresh: a.Refresh,
		Expire:  a.Expire,
	}
}

// Set data from configuration.AuthKeys
// implementation of AuthKeysRepository
func (a *AuthKeys) Set(c configuration.AuthKeys) {
	a.Access = c.Access
	a.Type = c.Type
	a.Refresh = c.Refresh
	a.Expire = c.Expire
}

// 147
