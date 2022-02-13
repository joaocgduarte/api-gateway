package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/plagioriginal/api-gateway/domain"
	"github.com/plagioriginal/api-gateway/helpers"
)

type AuthorizationMiddleware struct {
	tm domain.TokenManager
	uc domain.UsersClient
	ch domain.CookieHandler
	l  *log.Logger
}

// Returns a new instance of the middleware
func NewAuthorizationMiddleware(
	tm domain.TokenManager,
	uc domain.UsersClient,
	ch domain.CookieHandler,
	l *log.Logger,
) AuthorizationMiddleware {
	return AuthorizationMiddleware{
		tm: tm,
		uc: uc,
		ch: ch,
		l:  l,
	}
}

// Requires a valid token
func (aw AuthorizationMiddleware) RequireToken(allowedRoles []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString := aw.ch.GetAccessToken(r)

			if len(tokenString) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				helpers.JSON(w, r, "invalid token")
				return
			}

			token, err := aw.tm.ParseToken(tokenString)
			if err != nil || !aw.tm.IsTokenValid(token) {
				newTokens, err := aw.getNewTokens(r)
				if err != nil {
					aw.l.Printf("error fetching new tokens: %v\n", err)
					w.WriteHeader(http.StatusUnauthorized)
					helpers.JSON(w, r, "invalid token")
					return
				}

				aw.ch.GenerateCookiesFromTokens(w, newTokens.AccessToken, newTokens.RefreshToken)

				token, err = aw.tm.ParseToken(newTokens.AccessToken)
				if err != nil {
					aw.l.Printf("error parsing token: %v\n", err)
					w.WriteHeader(http.StatusUnauthorized)
					helpers.JSON(w, r, "invalid token")
					return
				}
			}

			userRole, err := aw.tm.GetTokenRole(token)
			if err != nil {
				aw.l.Printf("error fetching the role of the token: %v\n", err)
				w.WriteHeader(http.StatusUnauthorized)
				helpers.JSON(w, r, "invalid token")
				return
			}

			if len(allowedRoles) > 0 && !helpers.InArray(userRole, allowedRoles) {
				w.WriteHeader(http.StatusUnauthorized)
				helpers.JSON(w, r, "invalid token")
				return
			}
			userId, err := aw.tm.GetTokenIssuer(token)

			if err != nil {
				aw.l.Printf("error issuer of the token: %v\n", err)
				w.WriteHeader(http.StatusUnauthorized)
				helpers.JSON(w, r, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), "userId", userId)
			ctx = context.WithValue(ctx, "userRole", userRole)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

// Gets new tokens from the UsersClient
func (aw AuthorizationMiddleware) getNewTokens(r *http.Request) (domain.TokenResponse, error) {
	refreshToken := aw.ch.GetRefreshToken(r)
	aw.l.Println(refreshToken)

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	res, err := aw.uc.RefreshJWT(ctx, refreshToken)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	return domain.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
