package models

type AuthMode string

const (
	AuthModeBasic  AuthMode = "basic"  // Use basic auth for authentication with credentials defined in the config
	AuthModeBearer AuthMode = "bearer" // Use JWT for authentication with credentials defined in the config
	AuthModeOAuth  AuthMode = "oauth2" // Use OAuth2 for authentication
)

type Config struct {
	UserLogin         string
	UserSecret        string
	WApiUsername      string
	WApiPassword      string
	BaseUrl           string
	UseLogFile        bool
	JsonWebKey        string
	AuthMode          AuthMode
	OAuthIssuer       string
	OAuthAudience     string
	OAuthAllowedRoles []string // Use empty string to skip role checks, if defined, roles need to be separated by commas
}
