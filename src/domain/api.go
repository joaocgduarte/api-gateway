package domain

import "net/http"

type CookieHandler interface {
	GetAccessToken(r *http.Request) string
	GetRefreshToken(r *http.Request) string
	GenerateCookiesFromTokens(w http.ResponseWriter, accessToken string, refreshToken string)
}
