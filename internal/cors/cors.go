package cors

import (
	"escapade/internal/config"
	"fmt"

	"net/http"
	"strings"
)

// IsAllowed check can this site connect to server
func IsAllowed(origin string, origins []string) (allowed bool) {
	allowed = false
	for _, str := range origins {
		if str == origin {
			allowed = true
			break
		}
	}
	fmt.Println("Cant found:", origin)
	return
}

// SetCORS set cors headers
func SetCORS(rw http.ResponseWriter, cc config.CORSConfig, name string) {
	rw.Header().Set("Access-Control-Allow-Origin", name)
	rw.Header().Set("Access-Control-Allow-Headers", strings.Join(cc.Headers, ", "))
	rw.Header().Set("Access-Control-Allow-Credentials", cc.Credentials)
	rw.Header().Set("Access-Control-Allow-Methods", strings.Join(cc.Methods, ", "))
}

// getOrigin get domain connected to server
func GetOrigin(r *http.Request) string {
	return r.Header.Get("Origin")
}
