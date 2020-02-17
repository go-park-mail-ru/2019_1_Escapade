package oauth2server

import "gopkg.in/oauth2.v3"

const (
	ErrNoTokenConfiguration = "No token configuration given"
	ErrNoManager            = "No manager given"

	ErrUserNotFound

	AllowGetAccessRequest = true
)

var (
	AllowedResponseTypes = []oauth2.ResponseType{
		oauth2.Code,
		oauth2.Token,
	}
	AllowedGrantTypes = []oauth2.GrantType{
		oauth2.PasswordCredentials,
		oauth2.Refreshing,
	}
)
