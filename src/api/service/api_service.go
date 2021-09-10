package service

import (
	"context"
	"log"
	"time"

	"github.com/plagioriginal/api-gateway/domain"
	"github.com/plagioriginal/api-gateway/protos/protos"
)

type APIDefaultService struct {
	UsersClient    protos.UsersClient
	Logger         *log.Logger
	contextTimeout time.Duration
}

func New(uc protos.UsersClient, l *log.Logger, contextTimeout time.Duration) APIDefaultService {
	return APIDefaultService{
		UsersClient:    uc,
		Logger:         l,
		contextTimeout: contextTimeout,
	}
}

// Login route handler
func (as APIDefaultService) Login(ctx context.Context, loginRequest domain.LoginRequest) (*protos.TokenResponse, error) {
	_, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	adaptedRequest := &protos.LoginRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	}

	return as.UsersClient.Login(ctx, adaptedRequest)
}

// Refresh JWT token handler.
func (as APIDefaultService) RefreshJWT(ctx context.Context, refreshToken string) (*protos.TokenResponse, error) {
	_, cancel := context.WithTimeout(ctx, as.contextTimeout)
	defer cancel()

	refreshRequest := &protos.RefreshRequest{
		RefreshToken: refreshToken,
	}

	return as.UsersClient.Refresh(ctx, refreshRequest)
}

func (as APIDefaultService) AddUser(ctx context.Context, userRequest domain.AddUserRequest) (*protos.UserResponse, error) {
	return nil, nil
}

func (as APIDefaultService) ListProjects(ctx context.Context, listProjectReq domain.ListProjectsRequest) (interface{}, error) {
	return nil, nil
}
