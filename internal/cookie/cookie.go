package cookie

import (
	"escapade/internal/config"
	"escapade/internal/utils"
	"net/http"
)

// CreateID generate random string - session ID
func CreateID(length int) string {
	return utils.RandomString(length)
}

// CreateCookie create instance of cookie
func CreateCookie(value string, cc config.CookieConfig) (cookie *http.Cookie) {
	cookie = &http.Cookie{
		Name:     cc.NameCookie,
		Value:    value,
		Path:     cc.PathCookie,
		MaxAge:   cc.LifetimeCookie * 100000,
		HttpOnly: cc.HTTPOnly,
	}
	return
}

// GetSessionCookie get session cookie from request
func GetSessionCookie(r *http.Request, cc config.CookieConfig) (string, error) {
	session, err := r.Cookie(cc.NameCookie)
	if err != nil || session == nil || session.Value == "" {
		return "", err
	}
	return session.Value, err
}

// CreateAndSet creates cookie with value - value and sets it
func CreateAndSet(w http.ResponseWriter, cc config.CookieConfig, value string) {
	http.SetCookie(w, CreateCookie(value, cc))
}
