package middleware

import (
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
)

type CORS struct {
	c     config.CORS
	log   infrastructure.LoggerI
	trace infrastructure.ErrorTrace
}

func NewCORS(
	c config.CORS,
	log infrastructure.LoggerI,
	trace infrastructure.ErrorTrace,
) *CORS {
	return &CORS{
		c:     c,
		log:   log,
		trace: trace,
	}
}

func (c *CORS) Func(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			origin := getOrigin(r)
			c.log.Println("cors check")
			if !c.isAllowed(origin, c.c.Origins) {
				handlers.SendResult(
					rw,
					handlers.NewResult(
						http.StatusForbidden,
						nil,
						c.trace.New(ErrCors+origin),
					),
				)
				c.log.Println("cors no!!!!!!!!!!")
				return
			}
			c.log.Println("cors allow!")
			setCORS(rw, c.c, origin)
			next.ServeHTTP(rw, r)
			return
		},
	)
}

// IsAllowed check can this site connect to server
func (c *CORS) isAllowed(
	origin string,
	origins []string,
) (allowed bool) {
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
		c.log.Println("cors:", origin, "not allowed!")
	}
	return
}

// SetCORS set cors headers
func setCORS(
	rw http.ResponseWriter,
	cc config.CORS,
	name string,
) {
	rw.Header().Set(
		"Access-Control-Allow-Origin",
		name,
	)
	rw.Header().Set(
		"Access-Control-Allow-Headers",
		strings.Join(cc.Headers, ", "),
	)
	rw.Header().Set(
		"Access-Control-Allow-Credentials",
		cc.Credentials,
	)
	rw.Header().Set(
		"Access-Control-Allow-Methods",
		"GET, POST, DELETE, PUT, OPTIONS", // TODO а чего не из конфига?
	)
}

// GetOrigin get domain connected to server
func getOrigin(r *http.Request) string {
	return r.Header.Get("Origin")
}
