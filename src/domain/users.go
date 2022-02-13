package domain

import (
	"context"
	"net/http"

	users "github.com/plagioriginal/users-service-grpc/users"
)

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LoginRequest struct {
	Username string `validate:"required,min=3" json:"username"`
	Password string `validate:"required,min=8" json:"password"`
}

type AddUserRequest struct {
	Username string `validate:"required,min=3" json:"username"`
	Password string `validate:"required,min=3" json:"password"`
	Role     string `validate:"required" json:"role"`
	JwtToken string `validate:"required" json:"token"`
}

// Services to be used by the HTTP Handler.
// Should have all the gRPC clients into it's base, so it could do all the requests.
type UsersClient interface {
	Logout(ctx context.Context, refreshToken string) (*users.TokenResponse, error)
	Login(ctx context.Context, loginRequest LoginRequest) (*users.TokenResponse, error)
	RefreshJWT(ctx context.Context, refreshToken string) (*users.TokenResponse, error)
	AddUser(ctx context.Context, userRequest AddUserRequest) (*users.UserResponse, error)
}

type UsersHttpHandler interface {
	Logout(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshJWT(w http.ResponseWriter, r *http.Request)
	AddUser(w http.ResponseWriter, r *http.Request)
}
