package domain

import "github.com/golang-jwt/jwt"

// Manages token operations. Verifies if tokens are valid
// To be used on a middleware for the necessary http routes.
type TokenManager interface {
	// @todo: replace `interface{}` with proper JWT Object
	ParseToken(tokenString string) (*jwt.Token, error)
	IsTokenValid(token *jwt.Token) bool
	GetTokenRole(token *jwt.Token) (string, error)
	GetTokenIssuer(token *jwt.Token) (string, error)
}
