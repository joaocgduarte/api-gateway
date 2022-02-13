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

func New(
	uc users.UsersClient,
	l *log.Logger,
	contextTimeout time.Duration,
) domain.UsersClient {
	return GrpcUsersClient{
		UsersClient:    uc,
		Logger:         l,
		contextTimeout: contextTimeout,
	}
}

// Login route handler
func (as GrpcUsersClient) Login(ctx context.Context, loginRequest domain.LoginRequest) (*domain.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	res, err := as.UsersClient.Login(ctx, &users.LoginRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	})

	if err != nil {
		return nil, err
	}
	return mapGrpcTokenResponseToDomain(res), nil
}

// Refresh JWT token handler.
func (as GrpcUsersClient) RefreshJWT(ctx context.Context, refreshToken string) (*domain.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	res, err := as.UsersClient.Refresh(ctx, &users.RefreshRequest{
		RefreshToken: refreshToken,
	})

	if err != nil {
		return nil, err
	}
	return mapGrpcTokenResponseToDomain(res), nil
}

// Handles the user logout.
func (as GrpcUsersClient) Logout(ctx context.Context, refreshToken string) (*domain.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	res, err := as.UsersClient.Logout(ctx, &users.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}
	return mapGrpcTokenResponseToDomain(res), nil
}

func (as GrpcUsersClient) AddUser(ctx context.Context, userRequest domain.AddUserRequest) (*domain.User, error) {
	return nil, nil
}
