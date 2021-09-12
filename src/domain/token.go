package domain

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

// Manages token operations. Verifies if tokens are valid
// To be used on a middleware for the necessary http routes.
type TokenManager interface {
	// @todo: replace `interface{}` with proper JWT Object
	GetRefreshTokenFromCookies(r *http.Request) string
	GetJWTTokenFromCookies(r *http.Request) string
	ParseToken(tokenString string) (*jwt.Token, error)
	IsTokenValid(token *jwt.Token) bool
	GetTokenRole(token *jwt.Token) (string, error)
	GetTokenIssuer(token *jwt.Token) (string, error)
}
