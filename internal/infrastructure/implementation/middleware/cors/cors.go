package micors

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
)

// CORS is implementation of Middleware interface(package infrastructure)
// realize CORS(https://developer.mozilla.org/ru/docs/Web/HTTP/CORS)
type CORS struct {
	c     configuration.Cors
	log   infrastructure.Logger
	trace infrastructure.ErrorTrace
}

// New instance of CORS
func New(
	c configuration.CorsRepository,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
) (*CORS, error) {
	// check configuration repository given
	if c == nil {
		return nil, errors.New(ErrorNoConfiguration)
	}

	//overriding the nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}

	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}

	return &CORS{
		c:     c.Get(),
		log:   logger,
		trace: trace,
	}, nil
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
					c.log,
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
	cors configuration.Cors,
	name string,
) {
	rw.Header().Set(AllowOrigin, name)
	rw.Header().Set(AllowHeaders, strings.Join(cors.Headers, ", "))
	rw.Header().Set(AllowCredentials, cors.Credentials)
	rw.Header().Set(AllowMethods, strings.Join(cors.Methods, ", "))
}

// GetOrigin get domain connected to server
func getOrigin(r *http.Request) string {
	return r.Header.Get(Origin)
}
