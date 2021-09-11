package tokens

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/plagioriginal/api-gateway/domain"
)

// Our custom claimes for the JWT Token
type ClaimsWithRole struct {
	UserRoleSlug  string `json:"roleSlug"`
	UserRoleLabel string `json:"roleLabel"`
	Username      string `json:"username"`
	jwt.StandardClaims
}

// Object used to manage token/auth operations
type DefaultTokenManager struct {
	JWTSecret string
}

// Instantiates a new Token Manager
func NewTokenManager(jwtSecret string) domain.TokenManager {
	return DefaultTokenManager{
		JWTSecret: jwtSecret,
	}
}

// Gets the token from the bearer header.
func (t DefaultTokenManager) GetJWTTokenFromHeaders(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")

	if len(authHeader) == 0 {
		return ""
	}

	token := authHeader[len("Bearer "):]
	return token
}

// Gets the refresh token from the header.
func (t DefaultTokenManager) GetRefreshTokenFromHeaders(r *http.Request) string {
	return r.Header.Get("refresh-token")
}

// Gets the user ID from token
func (t DefaultTokenManager) GetTokenIssuer(token *jwt.Token) (string, error) {
	if claims, ok := token.Claims.(*ClaimsWithRole); ok && t.IsTokenValid(token) {
		return claims.StandardClaims.Issuer, nil
	}

	return "", domain.ErrInvalidToken
}

// Gets the user role from a jwt token
func (t DefaultTokenManager) GetTokenRole(token *jwt.Token) (string, error) {
	if claims, ok := token.Claims.(*ClaimsWithRole); ok && t.IsTokenValid(token) {
		return claims.UserRoleSlug, nil
	}

	return "", domain.ErrInvalidToken
}

// Checks if a token is valid
func (t DefaultTokenManager) IsTokenValid(token *jwt.Token) bool {
	return token.Valid
}

// Parses a JWT Token string to an object.
func (t DefaultTokenManager) ParseToken(tokenString string) (*jwt.Token, error) {
	key := []byte(t.JWTSecret)

	return jwt.ParseWithClaims(tokenString, &ClaimsWithRole{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
}
