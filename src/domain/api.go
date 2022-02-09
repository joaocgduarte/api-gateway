package domain

import (
	"context"

	users "github.com/plagioriginal/users-service-grpc/users"
)

type FailRestResponse struct {
	Errors string
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

type ListProjectsRequest struct {
	JwtToken string `validate:"required"`
	PageNum  uint
}

// Services to be used by the HTTP Handler.
// Should have all the gRPC clients into it's base, so it could do all the requests.
//
// @todo: define the `interfaces{}` with the Objects from the gRPC clients.
type APIService interface {
	Logout(ctx context.Context, refreshToken string) (*users.TokenResponse, error)
	Login(ctx context.Context, loginRequest LoginRequest) (*users.TokenResponse, error)
	RefreshJWT(ctx context.Context, refreshToken string) (*users.TokenResponse, error)
	AddUser(ctx context.Context, userRequest AddUserRequest) (*users.UserResponse, error)
	ListProjects(ctx context.Context, listProjectReq ListProjectsRequest) (interface{}, error)
}
