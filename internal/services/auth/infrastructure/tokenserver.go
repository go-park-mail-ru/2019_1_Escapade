package infrastructure

import (
	"net/http"

	"gopkg.in/oauth2.v3"
)

type TokenServer interface {
	ValidationBearerToken(r *http.Request) (ti oauth2.TokenInfo, err error)
	HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error
	HandleTokenRequest(w http.ResponseWriter, r *http.Request) error
}
