package models

// SessionToken represents OAUTH2 token
type SessionToken struct {
	AccessToken  string `json:"access_token" example:"123123123" `
	TokenType    string `json:"token_type"  example:"bearer" `
	ExpiresIn    int32  `json:"expires_in"  example:"86400" `
	RefreshToken string `json:"refresh_token" example:"321321321" `
}

// ErrorDescription represents OAUTH2 Error
type ErrorDescription struct {
	Error            string `json:"error" example:"unsupported_grant_type" `
	ErrorDescription string `json:"error_description" example:"The authorization grant type is not supported by the authorization server" `
}
