package auth

const (
	ErrNoHeaders        = "Token headers not found"
	ErrInvalidTokenType = "Wrong type of auth token"

	TimeFormat = "2006-01-02 15:04:05"

	HeaderAuthorization = "Authorization"
	HeaderAccess        = "Authorization-Access"
	HeaderType          = "Authorization-Type"
	HeaderRefresh       = "Authorization-Refresh"
	HeaderExpire        = "Authorization-Expire"

	CookieName = "access_token"

	AddrTest      = "%s/auth/test?access_token=%s"
	AddrDelete    = "%s/auth/delete?access_token=%s"
	AddrAuthorize = "/auth/authorize"
	AddrToken     = "/auth/token"
)
