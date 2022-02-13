package cookies

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/plagioriginal/api-gateway/domain"
)

const (
	accessTokenKey  string = "access-token"
	refreshTokenKey string = "refresh-token"
)

// Settings for the cookies.
type CookieSettings struct {
	Name  string
	Value string
}

type CookieHandler struct {
	cookieEncoder *securecookie.SecureCookie
}

func New(cookieEncoder *securecookie.SecureCookie) domain.CookieHandler {
	return CookieHandler{cookieEncoder: cookieEncoder}
}

func (c CookieHandler) GetAccessToken(r *http.Request) string {
	return c.getCookieValue(r, accessTokenKey)
}

func (c CookieHandler) GetRefreshToken(r *http.Request) string {
	return c.getCookieValue(r, refreshTokenKey)
}

// Gets a cookie value
func (c CookieHandler) getCookieValue(r *http.Request, cookieName string) string {
	if cookie, err := r.Cookie(cookieName); err == nil {
		var value string

		if err = c.cookieEncoder.Decode(cookieName, cookie.Value, &value); err == nil {
			return value
		}
	}

	return ""
}

// Generates the cookies from the Access and Refresh Tokens
func (c CookieHandler) GenerateCookiesFromTokens(w http.ResponseWriter, accessToken string, refreshToken string) {
	cookies := []CookieSettings{
		{
			Name:  accessTokenKey,
			Value: accessToken,
		},
		{
			Name:  refreshTokenKey,
			Value: refreshToken,
		},
	}

	for _, cookie := range cookies {
		c.setUpHttpOnlyCookie(w, cookie)
	}
}

// Sets up the http only cookies.
func (c CookieHandler) setUpHttpOnlyCookie(w http.ResponseWriter, cookie CookieSettings) {
	encodedValue, err := c.cookieEncoder.Encode(cookie.Name, cookie.Value)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookie.Name,
		Value:    encodedValue,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Secure:   true,
	})
}
