package users

import (
	"context"
	"log"
	"time"

	"github.com/plagioriginal/api-gateway/domain"
	users "github.com/plagioriginal/users-service-grpc/users"
)

type GrpcUsersClient struct {
	UsersClient    users.UsersClient
	Logger         *log.Logger
	contextTimeout time.Duration
}

func New(uc users.UsersClient, l *log.Logger, contextTimeout time.Duration) domain.UsersClient {
	return GrpcUsersClient{
		UsersClient:    uc,
		Logger:         l,
		contextTimeout: contextTimeout,
	}
}

// Login route handler
func (as GrpcUsersClient) Login(ctx context.Context, loginRequest domain.LoginRequest) (*users.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	return as.UsersClient.Login(ctx, &users.LoginRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	})
}

// Refresh JWT token handler.
func (as GrpcUsersClient) RefreshJWT(ctx context.Context, refreshToken string) (*users.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	return as.UsersClient.Refresh(ctx, &users.RefreshRequest{
		RefreshToken: refreshToken,
	})
}

// Handles the user logout.
func (as GrpcUsersClient) Logout(ctx context.Context, refreshToken string) (*users.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	return as.UsersClient.Logout(ctx, &users.RefreshRequest{
		RefreshToken: refreshToken,
	})
}

func (as GrpcUsersClient) AddUser(ctx context.Context, userRequest domain.AddUserRequest) (*users.UserResponse, error) {
	return nil, nil
}
