package cors

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"net/http"
	"strings"
)

// IsAllowed check can this site connect to server
func IsAllowed(origin string, origins []string) (allowed bool) {
	if origin == "" {
		return true
	}
	allowed = false
	for _, str := range origins {
		if str == origin {
			allowed = true
			break
		}
	}
	if !allowed {
		utils.Debug(false, "cors:", origin, "not allowed!")
	}
	return
}

// SetCORS set cors headers
func SetCORS(rw http.ResponseWriter, cc config.CORSConfig, name string) {
	rw.Header().Set("Access-Control-Allow-Origin", name)
	rw.Header().Set("Access-Control-Allow-Headers", strings.Join(cc.Headers, ", "))
	rw.Header().Set("Access-Control-Allow-Credentials", cc.Credentials)
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
}

// GetOrigin get domain connected to server
func GetOrigin(r *http.Request) string {
	return r.Header.Get("Origin")
}
