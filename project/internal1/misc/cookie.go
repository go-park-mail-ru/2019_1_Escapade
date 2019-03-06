package misc

import (
	"crypto/rand"
	"net/http"
	"time"
)

const (
	NameCookie     = "sessionid"
	LengthCookie   = 16
	years          = 0
	months         = 0
	days           = 7
	LifetimeCookie = days * 24 * 60
)

func CreateExpiration() time.Time {
	return time.Now().AddDate(years, months, days)
}

func CreateID() string {
	return randStr(LengthCookie)
}

func CreateCookie(value string) (cookie *http.Cookie) {
	cookie = &http.Cookie{}
	cookie.MaxAge = LifetimeCookie
	cookie.Name = NameCookie
	cookie.Value = value
	return
}

func GetSessionCookie(r *http.Request) string {
	session, err := r.Cookie(NameCookie)
	if err != nil || session == nil {
		return ""
	}
	return session.Value
}

func SetCookie(w http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(w, cookie)
}

func randStr(strSize int) string {

	var dictionary string

	dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
