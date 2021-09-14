package cookies

import (
	"net/http"
	"time"
)

type CookieSettings struct {
	Name  string
	Value string
}

type MultipleCookieSettings []CookieSettings

// Sets up the http only cookies.
func (cs CookieSettings) SetUpHttpOnly(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cs.Name,
		Value:    cs.Value,
		HttpOnly: true,
		Path:     "/",
		// SameSite: http.SameSiteNoneMode,
		Expires: time.Now().Add(time.Hour * 24 * 14),
		// Secure:   true,
	})
}

// Sets up http cookies for a series of cookies.
func (cs MultipleCookieSettings) SetUpHttpOnly(w http.ResponseWriter) {
	for _, cookie := range cs {
		cookie.SetUpHttpOnly(w)
	}
}

// Generates the cookies from the Access and Refresh Tokens
func GenerateCookiesFromTokens(w http.ResponseWriter, accessToken string, refreshToken string) {
	accessCookies := MultipleCookieSettings{
		{
			Name:  "access-token",
			Value: accessToken,
		},
		{
			Name:  "refresh-token",
			Value: refreshToken,
		},
	}

	accessCookies.SetUpHttpOnly(w)
}
