package domain

import (
	"context"
	"net/http"
)

type Role struct {
	Id        string `json:"id"`
	RoleLabel string `json:"role_label"`
	RoleSlug  string `json:"role_slug"`
}

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      Role   `json:"role"`
}

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
	Logout(ctx context.Context, refreshToken string) (*TokenResponse, error)
	Login(ctx context.Context, loginRequest LoginRequest) (*TokenResponse, error)
	RefreshJWT(ctx context.Context, refreshToken string) (*TokenResponse, error)
	AddUser(ctx context.Context, userRequest AddUserRequest) (*User, error)
}

type UsersHttpHandler interface {
	Logout(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshJWT(w http.ResponseWriter, r *http.Request)
	AddUser(w http.ResponseWriter, r *http.Request)
}
