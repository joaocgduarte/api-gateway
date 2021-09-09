package domain

import "context"

type LoginRequest struct {
	Username string
	Passwor  string
}

type AddUserRequest struct {
	Username string
	Password string
	Role     string
	JwtToken string
}

type ListProjectsRequest struct {
	JwtToken string
	PageNum  uint
}

// Services to be used by the HTTP Handler.
// Should have all the gRPC clients into it's base, so it could do all the requests.
//
// @todo: define the `interfaces{}` with the Objects from the gRPC clients.
type APIService interface {
	Login(ctx context.Context, loginRequest LoginRequest) (interface{}, error)
	RefreshJWT(ctx context.Context, refreshToken string) (interface{}, error)
	AddUser(ctx context.Context, userRequest AddUserRequest) (interface{}, error)
	ListProjects(ctx context.Context, listProjectReq ListProjectsRequest) (interface{}, error)
}
