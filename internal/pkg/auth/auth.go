package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/oauth2"
	"gopkg.in/oauth2.v3/models"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/cookie"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

func HashPassword(password, salt string) string {
	if password == "" {
		return password
	}

	hasher := sha256.New224()
	hasher.Write([]byte(password + salt))
	password = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	hasher.Write([]byte(salt + password))
	password = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return password
}

func update(rw http.ResponseWriter, token oauth2.Token,
	ca oauth2.Config) (oauth2.Token, error) {

	token.Expiry = time.Now()

	tokenSource := ca.TokenSource(context.Background(), &token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return oauth2.Token{}, err
	}
	return *newToken, err
}

func Check(rw http.ResponseWriter, r *http.Request,
	cc config.Cookie, ca config.Auth, client config.AuthClient) (string, error) {

	token, err := tokenFromHeadersNEW(r)
	if err != nil {
		err = nil //return "", err
	} else {
		return checkNEW(token, client)
	}
	// code below below will be deleted later
	///////////////////////////////////////////////////

	var (
		accessToken string
		updated     bool
	)

	// token given in cookie
	if !(cc.LifetimeHours == 0) {
		utils.Debug(false, "look in cookies")
		isReserve := false
		token, err = cookie.GetToken(r, cc, isReserve)
		if err != nil {
			isReserve = true
			token, err = cookie.GetToken(r, cc, isReserve)
		}
		if err == nil {
			utils.Debug(false, "all ok")
		} else {
			utils.Debug(false, "error catched", err.Error())
		}
		accessToken, token, updated, err = check(rw, r, false, token, ca, client)
		if err == nil {
			if updated {
				cookie.SetToken(rw, isReserve, token, cc)
			}
			return accessToken, err
		}
	}

	// token given in headers
	if accessToken == "" {
		token, err = tokenFromHeaders(r)
		if err != nil {
			return accessToken, err
		}
		accessToken, token, updated, err = check(rw, r, false, token, ca, client)
		if err == nil {
			if updated {
				setTokenToHeaders(rw, token)
			}
		}
	}
	return accessToken, err
}

func tokenFromHeadersNEW(r *http.Request) (oauth2.Token, error) {

	access := r.Header.Get("Authorization")
	if access == "" {
		return oauth2.Token{}, re.NoHeaders()
	}
	elements := strings.Split(access, " ")
	token := oauth2.Token{
		AccessToken: elements[1],
		TokenType:   elements[0],
	}
	return token, nil
}

// deprecated
func tokenFromHeaders(r *http.Request) (oauth2.Token, error) {

	token := oauth2.Token{
		AccessToken:  r.Header.Get("Authorization-Access"),
		TokenType:    r.Header.Get("Authorization-Type"),
		RefreshToken: r.Header.Get("Authorization-Refresh"),
	}
	if token.AccessToken == "" {
		return token, re.NoHeaders()
	}
	expireString := r.Header.Get("Authorization-Expire")
	token.Expiry, _ = time.Parse("2006-01-02 15:04:05", expireString)
	return token, nil
}

func setTokenToHeaders(rw http.ResponseWriter, token oauth2.Token) {

	rw.Header().Set("Authorization-Access", token.AccessToken)
	rw.Header().Set("Authorization-Type", token.TokenType)
	rw.Header().Set("Authorization-Refresh", token.RefreshToken)
	rw.Header().Set("Authorization-Expire", token.Expiry.Format("2006-01-02 15:04:05"))
	return
}

func checkNEW(token oauth2.Token, client config.AuthClient) (string, error) {

	resp, err := http.Get(fmt.Sprintf("%s/auth/test?access_token=%s",
		client.Address, token.AccessToken))
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

func check(rw http.ResponseWriter, r *http.Request, isReserve bool,
	token oauth2.Token, ca config.Auth, client config.AuthClient) (string, oauth2.Token, bool, error) {

	if token.TokenType != ca.TokenType {
		utils.Debug(false, "TokenType wrong! Get:", token.TokenType, ". Expected:", ca.TokenType)
		return "", oauth2.Token{}, false, re.ErrorTokenType()
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

	utils.Debug(false, "look at access ", token.AccessToken)
	utils.Debug(false, "look at type ", token.TokenType)
	utils.Debug(false, "look at expiry ", token.Expiry)
	utils.Debug(false, "look at refresh", token.RefreshToken)
	var updated bool
	resp, err := http.Get(fmt.Sprintf("%s/auth/test?access_token=%s",
		client.Address, token.AccessToken))
	if err != nil {
		token, err = update(rw, token, client.Config)
		resp, err = http.Get(fmt.Sprintf("%s/auth/test?access_token=%s",
			client.Address, token.AccessToken))
		updated = true
	}
	if err != nil {
		utils.Debug(false, "get cant sorry", err.Error())
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

func DeleteFromHeader(r *http.Request, client config.AuthClient) error {
	token, err := tokenFromHeaders(r)
	if err != nil {
		return err
	}

	_, err = http.Get(fmt.Sprintf("%s/delete?access_token=%s", client.Address, token.AccessToken))
	return err
}

func DeleteToken(rw http.ResponseWriter, r *http.Request,
	cc config.Cookie, client config.AuthClient) error {
	accessToken, err := cookie.GetCookie(r, "access_token")
	if err != nil {
		return http.ErrNoCookie
	}
	cookie.DeleteToken(rw, false, cc)
	cookie.DeleteToken(rw, true, cc)

	_, err = http.Get(fmt.Sprintf("%s/delete?access_token=%s", client.Address, accessToken))
	return err
}

func CreateTokenInCookies(rw http.ResponseWriter, name, password string,
	config oauth2.Config, cc config.Cookie) error {
	utils.Debug(false, "try create tokens")
	token, err := config.PasswordCredentialsToken(context.Background(), name, password)
	if err != nil {
		return err
	}
	cookie.SetToken(rw, false, *token, cc)

	token, err = config.PasswordCredentialsToken(context.Background(), name, password)
	if err != nil {
		return err
	}
	cookie.SetToken(rw, true, *token, cc)
	return err
}

func CreateTokenInHeaders(rw http.ResponseWriter, name, password string,
	config oauth2.Config) (*oauth2.Token, error) {
	utils.Debug(false, config.ClientID, config.ClientSecret, config.Endpoint, config.RedirectURL, config.Scopes)
	token, err := config.PasswordCredentialsToken(context.Background(), name, password)
	if err != nil {
		return token, err
	}
	setTokenToHeaders(rw, *token)
	return token, err
}
