package cookie

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
)

func Cookie(name string, value string, expire time.Time) (cookie *http.Cookie) {
	cookie = &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expire,
		HttpOnly: false,
	}
	return
}

// CreateCookie create instance of cookie
func CreateCookie(value string, cc config.SessionConfig) (cookie *http.Cookie) {
	cookie = &http.Cookie{
		Name:     cc.Name,
		Value:    value,
		Path:     cc.Path,
		MaxAge:   cc.LifetimeSeconds,
		HttpOnly: cc.HTTPOnly,
	}
	return
}

// GetSessionCookie get session cookie from request
func GetSessionCookie1(r *http.Request, cc config.SessionConfig) (string, error) {
	session, err := r.Cookie(cc.Name)
	if err != nil || session == nil || session.Value == "" {
		return "", err
	}
	return session.Value, err
}

// CreateAndSet creates cookie with value - value and sets it
func CreateAndSet(w http.ResponseWriter, cc config.SessionConfig, value string) {
	http.SetCookie(w, CreateCookie(value, cc))
}
