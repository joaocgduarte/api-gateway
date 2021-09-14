package cookies

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
)

var (
	hashSecret    = []byte(os.Getenv("COOKIE_ENCODE_SECRET"))
	cookieEncoder = securecookie.New(hashSecret, nil)

	accessTokenKey  = "access-token"
	refreshTokenKey = "refresh-token"
)

// Settings for the cookies.
type CookieSettings struct {
	Name  string
	Value string
}

func GetAccessToken(r *http.Request) string {
	return GetCookieValue(r, accessTokenKey)
}

func GetRefreshToken(r *http.Request) string {
	return GetCookieValue(r, refreshTokenKey)
}

// Gets a cookie value
func GetCookieValue(r *http.Request, cookieName string) string {
	if cookie, err := r.Cookie(cookieName); err == nil {
		var value string

		if err = cookieEncoder.Decode(cookieName, cookie.Value, &value); err == nil {
			return value
		}
	}

	return ""
}

// Encodes a cookie value.
func EncodeCookieValue(cs CookieSettings) (string, error) {
	return cookieEncoder.Encode(cs.Name, cs.Value)
}

// Sets up the http only cookies.
func (cs CookieSettings) SetUpHttpOnly(w http.ResponseWriter) {
	encodedValue, err := EncodeCookieValue(cs)

	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cs.Name,
		Value:    encodedValue,
		HttpOnly: true,
		Path:     "/",
		// SameSite: http.SameSiteNoneMode,
		Expires: time.Now().Add(time.Hour * 24 * 14),
		// Secure:   true,
	})
}

type MultipleCookieSettings []CookieSettings

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
			Name:  accessTokenKey,
			Value: accessToken,
		},
		{
			Name:  refreshTokenKey,
			Value: refreshToken,
		},
	}

	accessCookies.SetUpHttpOnly(w)
}
