package misc

import (
	"crypto/rand"
	"net/http"
	"time"
)

const (
	NameCookie      = "sessionid"
	pathCoolie      = "/"
	LengthCookie    = 16
	years           = 0
	months          = 0
	days            = 7
	LifetimeCookie  = days * 24 * 60
	LengthImageName = 24
)

func CreateExpiration() time.Time {
	return time.Now().AddDate(years, months, days)
}

func CreateID() string {
	return randStr(LengthCookie)
}

func CreateImageName() string {
	return randStr(LengthImageName)
}

func CreateCookie(value string) (cookie *http.Cookie) {
	cookie = &http.Cookie{
		Name:     NameCookie,
		Value:    value,
		Path:     pathCoolie,
		MaxAge:   LifetimeCookie,
		HttpOnly: true,
	}

	return
}

func GetSessionCookie(r *http.Request) (string, error) {
	session, err := r.Cookie(NameCookie)
	if err != nil || session == nil {
		return "", err
	}
	return session.Value, err
}

func SetCookie(w http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(w, cookie)
}

// CreateAndSet creates cookie with value - value and sets it
func CreateAndSet(w http.ResponseWriter, value string) {
	http.SetCookie(w, CreateCookie(value))
}

func randStr(strSize int) string {

	dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
