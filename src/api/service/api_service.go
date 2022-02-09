package service

import (
	"context"
	"log"
	"time"

	"github.com/plagioriginal/api-gateway/domain"
	users "github.com/plagioriginal/users-service-grpc/users"
)

type APIDefaultService struct {
	UsersClient    users.UsersClient
	Logger         *log.Logger
	contextTimeout time.Duration
}

func New(uc users.UsersClient, l *log.Logger, contextTimeout time.Duration) domain.APIService {
	return APIDefaultService{
		UsersClient:    uc,
		Logger:         l,
		contextTimeout: contextTimeout,
	}
}

// Login route handler
func (as APIDefaultService) Login(ctx context.Context, loginRequest domain.LoginRequest) (*users.TokenResponse, error) {
	_, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	adaptedRequest := &users.LoginRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	}

	return as.UsersClient.Login(ctx, adaptedRequest)
}

// Refresh JWT token handler.
func (as APIDefaultService) RefreshJWT(ctx context.Context, refreshToken string) (*users.TokenResponse, error) {
	_, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	refreshRequest := &users.RefreshRequest{
		RefreshToken: refreshToken,
	}

	return as.UsersClient.Refresh(ctx, refreshRequest)
}

// Handles the user logout.
func (as APIDefaultService) Logout(ctx context.Context, refreshToken string) (*users.TokenResponse, error) {
	_, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	refreshRequest := &users.RefreshRequest{
		RefreshToken: refreshToken,
	}

	return as.UsersClient.Logout(ctx, refreshRequest)
}

func (as APIDefaultService) AddUser(ctx context.Context, userRequest domain.AddUserRequest) (*users.UserResponse, error) {
	return nil, nil
}

func (as APIDefaultService) ListProjects(ctx context.Context, listProjectReq domain.ListProjectsRequest) (interface{}, error) {
	return nil, nil
}
