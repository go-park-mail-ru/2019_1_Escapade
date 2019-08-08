package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/cookie"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"golang.org/x/oauth2"
	"gopkg.in/oauth2.v3/models"
)

var AuthServerURL = "http://localhost:3003/auth"

func update(rw http.ResponseWriter, r *http.Request, config oauth2.Config,
	refreshToken, accessTokenKey, tokenTypeKey, refreshTokenKey,
	expireKey string) (string, error) {

	saved := &oauth2.Token{}
	saved.RefreshToken = refreshToken
	saved.Expiry = time.Now()

	tokenSource := config.TokenSource(context.Background(), saved)
	newToken, err := tokenSource.Token()
	if err != nil {
		utils.Debug(false, "error while updating", err.Error())
		return "", err
	}

	utils.Debug(false, "set old", saved.AccessToken, saved.Expiry)
	utils.Debug(false, "set new", newToken.AccessToken, newToken.Expiry)

	cookieExpire := time.Now().Add(time.Hour * 24 * 30 * 3)
	http.SetCookie(rw, cookie.Cookie(accessTokenKey, newToken.AccessToken, cookieExpire))
	http.SetCookie(rw, cookie.Cookie(tokenTypeKey, newToken.TokenType, cookieExpire))
	http.SetCookie(rw, cookie.Cookie(refreshTokenKey, newToken.RefreshToken, cookieExpire))
	http.SetCookie(rw, cookie.Cookie(expireKey, newToken.Expiry.Format("2006-01-02 15:04:05"), cookieExpire))

	return newToken.AccessToken, err
}

func Check(rw http.ResponseWriter, r *http.Request, config oauth2.Config) (string, error) {
	accessToken, err := check(rw, r, config, "access_token", "token_type", "refresh_token", "expire")
	if err != nil {
		accessToken, err = check(rw, r, config, "r_access_token", "r_token_type", "r_refresh_token", "r_expire")
	}
	return accessToken, err
}

func check(rw http.ResponseWriter, r *http.Request, config oauth2.Config, accessTokenKey, tokenTypeKey,
	refreshTokenKey, expireStringKey string) (string, error) {

	var (
		accessToken, tokenType, refreshToken, expireString, newAccessToken string
		expire                                                             time.Time
		err, err1, err2, err3, err4                                        error
	)

	accessToken, err1 = cookie.GetCookie(r, accessTokenKey)
	tokenType, err2 = cookie.GetCookie(r, tokenTypeKey)
	refreshToken, err3 = cookie.GetCookie(r, refreshTokenKey)
	expireString, err4 = cookie.GetCookie(r, expireStringKey)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return "", err
	}

	if tokenType != "Bearer" {
		utils.Debug(false, "tokenType:", tokenType)
		return "", re.ErrorTokenType()
	}

	expire, err = time.Parse("2006-01-02 15:04:05", expireString)
	if err != nil {
		return "", err
	}

	now, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return "", err
	}

	if expire.Before(now) {
		newAccessToken, err = update(rw, r, config,
			refreshToken, accessTokenKey, refreshTokenKey,
			refreshTokenKey, expireStringKey)
		if err != nil {
			return "", err
		}
	} else {
		newAccessToken = accessToken
	}

	resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s",
		AuthServerURL, newAccessToken))
	if err != nil {
		utils.Debug(false, "get cant sorry", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	token := models.Token{}

	err = json.NewDecoder(resp.Body).Decode(&token)

	if err != nil {
		return "", err
	} else {
		utils.Debug(false, "!!!!!!!!!!!!!!!!!!!!!token:", token.GetClientID(), token.GetCode(), token.GetScope(), token.GetAccess())
	}
	return token.GetUserID(), err
}

func DeleteToken(rw http.ResponseWriter, r *http.Request) error {
	cook, err := r.Cookie("access_token")

	http.SetCookie(rw, cookie.Cookie("access_token", "", time.Unix(0, 0)))
	http.SetCookie(rw, cookie.Cookie("token_type", "", time.Unix(0, 0)))
	http.SetCookie(rw, cookie.Cookie("refresh_token", "", time.Unix(0, 0)))
	http.SetCookie(rw, cookie.Cookie("expire", "", time.Unix(0, 0)))

	http.SetCookie(rw, cookie.Cookie("r_access_token", "", time.Unix(0, 0)))
	http.SetCookie(rw, cookie.Cookie("r_token_type", "", time.Unix(0, 0)))
	http.SetCookie(rw, cookie.Cookie("r_refresh_token", "", time.Unix(0, 0)))
	http.SetCookie(rw, cookie.Cookie("r_expire", "", time.Unix(0, 0)))

	if err != nil || cook == nil || cook.Value == "" {
		return http.ErrNoCookie
	}

	_, err = http.Get(fmt.Sprintf("%s/delete?access_token=%s", AuthServerURL, cook.Value))
	return err
}

func CreateToken(rw http.ResponseWriter, config oauth2.Config,
	name, password string) error {
	token, err := config.PasswordCredentialsToken(context.Background(), name, password)
	if err != nil {
		return err
	}
	cookieExpire := time.Now().Add(time.Hour * 24 * 30 * 3)
	http.SetCookie(rw, cookie.Cookie("access_token", token.AccessToken, cookieExpire))
	http.SetCookie(rw, cookie.Cookie("token_type", token.TokenType, cookieExpire))
	http.SetCookie(rw, cookie.Cookie("refresh_token", token.RefreshToken, cookieExpire))
	http.SetCookie(rw, cookie.Cookie("expire", token.Expiry.Format("2006-01-02 15:04:05"), cookieExpire))

	token, err = config.PasswordCredentialsToken(context.Background(), name, password)
	if err != nil {
		return err
	}
	http.SetCookie(rw, cookie.Cookie("r_access_token", token.AccessToken, cookieExpire))
	http.SetCookie(rw, cookie.Cookie("r_token_type", token.TokenType, cookieExpire))
	http.SetCookie(rw, cookie.Cookie("r_refresh_token", token.RefreshToken, cookieExpire))
	http.SetCookie(rw, cookie.Cookie("r_expire", token.Expiry.Format("2006-01-02 15:04:05"), cookieExpire))
	return err
}
