package domain

// Manages token operations. Verifies if tokens are valid
// To be used on a middleware for the necessary http routes.
type TokenManager interface {
	// @todo: replace `interface{}` with proper JWT Object
	ParseToken(token string) (interface{}, error)
	IsTokenValid(token string) bool
	GetTokenRole(token string) (string, error)
	GetTokenIssuer(token interface{}) (string, error)
}
