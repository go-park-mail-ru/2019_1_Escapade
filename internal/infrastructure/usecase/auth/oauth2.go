package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/oauth2.v3/models"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
)

type OAuth2 struct {
	entity.Auth
	Config oauth2.Config

	trace infrastructure.ErrorTrace
	log   infrastructure.LoggerI
}

func NewOAuth2(
	rep infrastructure.AuthServiceRepositoryI,
	trace infrastructure.ErrorTrace,
	log infrastructure.LoggerI,
) *OAuth2 {
	authFromRep := rep.Get()
	return &OAuth2{
		Auth: authFromRep,
		Config: oauth2.Config{
			ClientID:     authFromRep.Client.ClientID,
			ClientSecret: authFromRep.Client.ClientSecret,
			Scopes:       authFromRep.Client.Scopes,
			RedirectURL:  authFromRep.Client.RedirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:  authFromRep.Client.Address + AddrAuthorize,
				TokenURL: authFromRep.Client.Address + AddrToken,
			}},
		trace: trace,
		log:   log,
	}
}

func (service *OAuth2) HashPassword(password string) string {
	if password == "" {
		return password
	}

	hasher := sha256.New224()
	hasher.Write([]byte(password + service.Salt))
	password = base64.URLEncoding.EncodeToString(
		hasher.Sum(nil),
	)
	hasher.Write([]byte(service.Salt + password))
	password = base64.URLEncoding.EncodeToString(
		hasher.Sum(nil),
	)
	return password
}

func (service *OAuth2) update(
	rw http.ResponseWriter,
	token oauth2.Token,
) (oauth2.Token, error) {

	token.Expiry = time.Now()

	tokenSource := service.Config.TokenSource(
		context.Background(),
		&token,
	)
	newToken, err := tokenSource.Token()
	if err != nil {
		return oauth2.Token{}, err
	}
	return *newToken, err
}

func (service *OAuth2) Check(
	rw http.ResponseWriter,
	r *http.Request,
) (string, error) {
	token, err := service.tokenFromHeadersNEW(r)
	if err != nil {
		err = nil //return "", err
	} else {
		return service.checkNEW(token)
	}
	// ! code below below will be deleted later
	///////////////////////////////////////////////////

	var (
		accessToken string
		updated     bool
	)

	// token given in cookie
	if !(service.Cookie.LifetimeHours == 0) {
		service.log.Println("look in cookies")
		isReserve := false
		token, err = GetToken(r, service.Cookie, isReserve)
		if err != nil {
			isReserve = true
			token, err = GetToken(r, service.Cookie, isReserve)
		}
		if err == nil {
			service.log.Println("all ok")
		} else {
			service.log.Println("error catched", err.Error())
		}
		accessToken, token, updated, err = service.check(
			rw,
			r,
			false,
			token,
		)
		if err == nil {
			if updated {
				SetToken(rw, isReserve, token, service.Cookie)
			}
			return accessToken, err
		}
	}

	// token given in headers
	if accessToken == "" {
		token, err = service.tokenFromHeadersNEW(r)
		if err != nil {
			return accessToken, err
		}
		accessToken, token, updated, err = service.check(
			rw,
			r,
			false,
			token,
		)
		if err == nil {
			if updated {
				service.setTokenToHeaders(rw, token)
			}
		}
	}
	return accessToken, err
}

func (service *OAuth2) tokenFromHeadersNEW(
	r *http.Request,
) (oauth2.Token, error) {
	access := r.Header.Get(HeaderAuthorization)
	if access == "" {
		return oauth2.Token{}, service.trace.New(ErrNoHeaders)
	}
	elements := strings.Split(access, " ")
	token := oauth2.Token{
		AccessToken: elements[1],
		TokenType:   elements[0],
	}
	return token, nil
}

// ! Deprecated
func (service *OAuth2) tokenFromHeaders(
	r *http.Request,
) (oauth2.Token, error) {
	token := oauth2.Token{
		AccessToken:  r.Header.Get(HeaderAccess),
		TokenType:    r.Header.Get(HeaderType),
		RefreshToken: r.Header.Get(HeaderRefresh),
	}
	if token.AccessToken == "" {
		return token, service.trace.New(ErrNoHeaders)
	}
	expireString := r.Header.Get(HeaderExpire)
	token.Expiry, _ = time.Parse(TimeFormat, expireString)
	return token, nil
}

func (service *OAuth2) setTokenToHeaders(
	rw http.ResponseWriter,
	token oauth2.Token,
) {
	rw.Header().Set(HeaderAccess, token.AccessToken)
	rw.Header().Set(HeaderType, token.TokenType)
	rw.Header().Set(HeaderRefresh, token.RefreshToken)
	rw.Header().Set(HeaderExpire, token.Expiry.Format(TimeFormat))
	return
}

func (service *OAuth2) checkNEW(
	token oauth2.Token,
) (string, error) {
	resp, err := http.Get(fmt.Sprintf(
		AddrTest,
		service.Client.Address,
		token.AccessToken,
	))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokenModel := models.Token{}

	// we dont use json.Decoder cause https://ahmet.im/blog/golang-json-decoder-pitfalls/
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(bytes, &tokenModel)
	if err != nil {
		return "", err
	}
	return tokenModel.GetUserID(), err
}

// ! Deprecated
func (service *OAuth2) check(
	rw http.ResponseWriter,
	r *http.Request,
	sReserve bool,
	token oauth2.Token,
) (string, oauth2.Token, bool, error) {
	if token.TokenType != service.Auth.Auth.TokenType {
		service.log.Println("TokenType wrong! Get:", token.TokenType, ". Expected:", service.Auth.Auth.TokenType)
		return "",
			oauth2.Token{},
			false,
			service.trace.New(ErrInvalidTokenType)
	}

	// now, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	// if err != nil {
	// 	return "", oauth2.Token{}, false, err
	// }

	// var updated bool
	// if token.Expiry.Before(now) {
	// 	updated = true
	// 	utils.Debug(false, "before, go to updare")
	// 	token, err = update(rw, token, client.Config)
	// 	if err != nil {
	// 		return "", token, updated, err
	// 	}
	// }
	//accessToken = token.AccessToken

	service.log.Println("look at access ", token.AccessToken)
	service.log.Println("look at type ", token.TokenType)
	service.log.Println("look at expiry ", token.Expiry)
	service.log.Println("look at refresh", token.RefreshToken)
	var updated bool
	resp, err := http.Get(fmt.Sprintf(
		AddrTest,
		service.Client.Address,
		token.AccessToken,
	))
	if err != nil {
		token, err = service.update(rw, token)
		resp, err = http.Get(fmt.Sprintf(
			AddrTest,
			service.Client.Address,
			token.AccessToken,
		))
		updated = true
	}
	if err != nil {
		service.log.Println("get cant sorry", err.Error())
		return "", token, updated, err
	}
	defer resp.Body.Close()

	tokenModel := models.Token{}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", token, updated, err
	}
	err = json.Unmarshal(bytes, &tokenModel)
	if err != nil {
		return "", token, updated, err
	}
	return tokenModel.GetUserID(), token, updated, err
}

func (service *OAuth2) DeleteFromHeader(
	r *http.Request,
) error {
	token, err := service.tokenFromHeaders(r)
	if err != nil {
		return err
	}

	_, err = http.Get(fmt.Sprintf(
		AddrDelete,
		service.Client.Address,
		token.AccessToken,
	))
	return err
}

// ! Deprecated
func (service *OAuth2) DeleteToken(
	rw http.ResponseWriter,
	r *http.Request,
) error {
	accessToken, err := GetCookie(r, CookieName)
	if err != nil {
		return http.ErrNoCookie
	}
	DeleteToken(rw, false, service.Cookie)
	DeleteToken(rw, true, service.Cookie)

	_, err = http.Get(fmt.Sprintf(
		AddrDelete,
		service.Client.Address,
		accessToken,
	))
	return err
}

func (service *OAuth2) CreateToken(
	rw http.ResponseWriter,
	name, password string,
) error {
	_, err := service.CreateTokenInHeaders(rw, name, password)
	if err != nil {
		return err
	}

	err = service.CreateTokenInCookies(rw, name, password)
	if err != nil {
		return err
	}

	return nil
}

// ! Deprecated
func (service *OAuth2) CreateTokenInCookies(
	rw http.ResponseWriter,
	name, password string,
) error {
	service.log.Println("try create tokens")
	token, err := service.Config.PasswordCredentialsToken(
		context.Background(),
		name,
		password,
	)
	if err != nil {
		return err
	}
	SetToken(rw, false, *token, service.Cookie)

	token, err = service.Config.PasswordCredentialsToken(
		context.Background(),
		name,
		password,
	)
	if err != nil {
		return err
	}
	SetToken(rw, true, *token, service.Cookie)
	return err
}

func (service *OAuth2) CreateTokenInHeaders(
	rw http.ResponseWriter,
	name, password string,
) (*oauth2.Token, error) {
	service.log.Println(false, service.Config.ClientID,
		service.Config.ClientSecret, service.Config.Endpoint,
		service.Config.RedirectURL, service.Config.Scopes)
	token, err := service.Config.PasswordCredentialsToken(context.Background(), name, password)
	if err != nil {
		return token, err
	}
	service.setTokenToHeaders(rw, *token)
	return token, err
}

func SetCookie(
	w http.ResponseWriter,
	name, value string,
	cc config.Cookie,
) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     cc.Path,
		Expires:  time.Now().Add(time.Hour * time.Duration(cc.LifetimeHours)),
		HttpOnly: cc.HTTPOnly,
	}
	http.SetCookie(w, cookie)
	return
}

func DeleteCookie(
	w http.ResponseWriter,
	name string,
	cc config.Cookie,
) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   "",
		Path:    cc.Path,
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
	return
}

func GetCookie(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(key)
	if err != nil || cookie == nil {
		return "", http.ErrNoCookie
	}
	return cookie.Value, nil
}

func GetToken(
	r *http.Request,
	cc config.Cookie,
	isReserve bool,
) (oauth2.Token, error) {
	var (
		token        oauth2.Token
		expireString string
		err          error
		aKey         = cc.Auth.AccessToken
		tKey         = cc.Auth.TokenType
		rKey         = cc.Auth.RefreshToken
		eKey         = cc.Auth.Expire
	)
	if isReserve {
		aKey = cc.Auth.ReservePrefix + aKey
		tKey = cc.Auth.ReservePrefix + tKey
		rKey = cc.Auth.ReservePrefix + rKey
		eKey = cc.Auth.ReservePrefix + eKey
	}
	token.AccessToken, err = GetCookie(r, aKey)
	if err != nil {
		return token, err
	}
	token.TokenType, err = GetCookie(r, tKey)
	if err != nil {
		return token, err
	}
	token.RefreshToken, err = GetCookie(r, rKey)
	if err != nil {
		return token, err
	}
	expireString, err = GetCookie(r, eKey)
	if err != nil {
		return token, err
	}
	token.Expiry, err = time.Parse(TimeFormat, expireString)

	return token, err
}

func SetToken(
	rw http.ResponseWriter,
	isReserve bool,
	token oauth2.Token,
	cc config.Cookie,
) {
	var (
		aKey = cc.Auth.AccessToken
		tKey = cc.Auth.TokenType
		rKey = cc.Auth.RefreshToken
		eKey = cc.Auth.Expire
	)
	if isReserve {
		aKey = cc.Auth.ReservePrefix + aKey
		tKey = cc.Auth.ReservePrefix + tKey
		rKey = cc.Auth.ReservePrefix + rKey
		eKey = cc.Auth.ReservePrefix + eKey
	}

	SetCookie(rw, aKey, token.AccessToken, cc)
	SetCookie(rw, tKey, token.TokenType, cc)
	SetCookie(rw, rKey, token.RefreshToken, cc)
	SetCookie(rw, eKey, token.Expiry.Format(TimeFormat), cc)
}

func DeleteToken(
	rw http.ResponseWriter,
	isReserve bool,
	cc config.Cookie,
) {
	var (
		aKey = cc.Auth.TokenType
		tKey = cc.Auth.TokenType
		rKey = cc.Auth.RefreshToken
		eKey = cc.Auth.Expire
	)
	if isReserve {
		aKey = cc.Auth.ReservePrefix + aKey
		tKey = cc.Auth.ReservePrefix + tKey
		rKey = cc.Auth.ReservePrefix + rKey
		eKey = cc.Auth.ReservePrefix + eKey
	}

	DeleteCookie(rw, aKey, cc)
	DeleteCookie(rw, tKey, cc)
	DeleteCookie(rw, rKey, cc)
	DeleteCookie(rw, eKey, cc)
}
