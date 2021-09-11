package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/plagioriginal/api-gateway/domain"
	"github.com/plagioriginal/api-gateway/http_renderer"
	"github.com/plagioriginal/api-gateway/protos/protos"
)

type AuthorizationMiddleware struct {
	tm domain.TokenManager
	uc protos.UsersClient
	l  *log.Logger
}

// Returns a new instance of the middleware
func NewAuthorizationMiddleware(tm domain.TokenManager, uc protos.UsersClient, l *log.Logger) AuthorizationMiddleware {
	return AuthorizationMiddleware{
		tm: tm,
		uc: uc,
		l:  l,
	}
}

// Requires a valid token
func (aw AuthorizationMiddleware) RequireValidToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenString := aw.tm.GetJWTTokenFromHeaders(r)

		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			http_renderer.JSON(w, r, domain.FailRestResponse{Errors: domain.ErrInvalidToken.Error()})
			return
		}

		ctx := r.Context()

		token, err := aw.tm.ParseToken(tokenString)

		if err != nil || !aw.tm.IsTokenValid(token) {
			aw.l.Println("requesting new tokens...")
			newTokens, err := aw.getNewTokens(r)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				http_renderer.JSON(w, r, domain.FailRestResponse{Errors: domain.ErrInvalidToken.Error()})
				return
			}

			ctx = context.WithValue(ctx, "newJwtToken", newTokens.AccessToken)
			ctx = context.WithValue(ctx, "newRefreshToken", newTokens.RefreshToken)
			token, _ = aw.tm.ParseToken(newTokens.AccessToken)
		}

		userId, _ := aw.tm.GetTokenIssuer(token)
		userRole, _ := aw.tm.GetTokenRole(token)

		ctx = context.WithValue(ctx, "userId", userId)
		ctx = context.WithValue(ctx, "userRole", userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Requires a valid token from an Admin role
func (aw AuthorizationMiddleware) RequireAdminValidToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenString := aw.tm.GetJWTTokenFromHeaders(r)

		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			http_renderer.JSON(w, r, domain.FailRestResponse{Errors: domain.ErrInvalidToken.Error()})
			return
		}

		ctx := r.Context()

		token, err := aw.tm.ParseToken(tokenString)

		if err != nil || !aw.tm.IsTokenValid(token) {
			aw.l.Println("requesting new tokens...")
			newTokens, err := aw.getNewTokens(r)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				http_renderer.JSON(w, r, domain.FailRestResponse{Errors: domain.ErrInvalidToken.Error()})
				return
			}

			ctx = context.WithValue(ctx, "newAccessTokens", newTokens)
			token, _ = aw.tm.ParseToken(newTokens.AccessToken)
		}

		userRole, _ := aw.tm.GetTokenRole(token)

		if userRole != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			http_renderer.JSON(w, r, domain.FailRestResponse{Errors: domain.ErrInvalidToken.Error()})
			return
		}

		userId, _ := aw.tm.GetTokenIssuer(token)

		ctx = context.WithValue(ctx, "userId", userId)
		ctx = context.WithValue(ctx, "userRole", userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Gets new tokens from the UsersClient
func (aw AuthorizationMiddleware) getNewTokens(r *http.Request) (*protos.TokenResponse, error) {
	refreshToken := aw.tm.GetRefreshTokenFromHeaders(r)

	return aw.uc.Refresh(r.Context(), &protos.RefreshRequest{
		RefreshToken: refreshToken,
	})
}
