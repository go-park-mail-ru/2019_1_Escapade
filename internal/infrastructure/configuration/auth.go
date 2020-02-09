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
	Salt            string
	Cookie          Cookie
	TokenGeneration TokenGeneration
	Client          AuthClient
}

// func NewConfiguration(
// 	cookie Cookie,
// 	auth Auth,
// 	client AuthClient,
// ) Configuration {
// 	return Configuration{
// 		Cookie: cookie,
// 		Auth:   auth,
// 		Client: client,
// 	}
// }

type TokenGenerationRepository interface {
	Get() TokenGeneration
	Set(TokenGeneration)
}

type TokenGeneration struct {
	AccessExpire, RefreshExpire time.Duration
	IsGenerateRefresh           bool
	TokenType                   string
}

// func NewAuth(
// 	salt string,
// 	accessTokenExpire, refreshTokenExpire time.Duration,
// 	isGenerateRefresh, withReserve bool,
// 	tokenType string,
// 	whiteList []AuthClient,
// ) Auth {
// 	return Auth{
// 		Salt:               salt,
// 		AccessTokenExpire:  accessTokenExpire,
// 		RefreshTokenExpire: refreshTokenExpire,
// 		IsGenerateRefresh:  isGenerateRefresh,
// 		WithReserve:        withReserve,
// 		TokenType:          tokenType,
// 		WhiteList:          whiteList,
// 	}
// }

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

// func NewAuthClient(
// 	address, clientID, clientSecret string,
// 	scopes []string,
// 	redirectURL string,
// ) AuthClient {
// 	return AuthClient{
// 		Address:      address,
// 		ClientID:     clientID,
// 		ClientSecret: clientSecret,
// 		Scopes:       scopes,
// 		RedirectURL:  redirectURL,
// 	}
// }

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

// func NewCookie(
// 	path string,
// 	lifetimeHours int,
// 	HTTPOnly bool,
// 	auth AuthCookie,
// ) Cookie {
// 	return Cookie{
// 		Path:          path,
// 		LifetimeHours: lifetimeHours,
// 		HTTPOnly:      HTTPOnly,
// 		Auth:          auth,
// 	}
// }

// func NewAuthCookie(
// 	accessToken, tokenType, refreshToken string,
// 	expire, reservePrefix string,
// ) AuthCookie {
// 	return AuthCookie{
// 		AccessToken:   accessToken,
// 		TokenType:     tokenType,
// 		RefreshToken:  refreshToken,
// 		Expire:        expire,
// 		ReservePrefix: reservePrefix,
// 	}
// }
