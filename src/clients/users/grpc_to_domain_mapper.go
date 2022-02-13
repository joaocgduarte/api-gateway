package users

import (
	"github.com/plagioriginal/api-gateway/domain"
	usersGrpc "github.com/plagioriginal/users-service-grpc/users"
)

func mapGrpcTokenResponseToDomain(in *usersGrpc.TokenResponse) *domain.TokenResponse {
	if in == nil {
		return nil
	}

	user := domain.User{}
	if in.GetUser() != nil {
		user = *mapGrpcUserResponseToDomain(in.GetUser())
	}

	return &domain.TokenResponse{
		AccessToken:  in.GetAccessToken(),
		RefreshToken: in.GetRefreshToken(),
		User:         user,
	}
}

func mapGrpcUserResponseToDomain(in *usersGrpc.UserResponse) *domain.User {
	if in == nil {
		return nil
	}

	res := &domain.User{
		Id:        in.GetId(),
		Username:  in.GetUsername(),
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Role:      domain.Role{},
	}
	if in.Role != nil {
		res.Role = domain.Role{
			Id:        in.GetRole().GetId(),
			RoleLabel: in.GetRole().GetRoleLabel(),
			RoleSlug:  in.GetRole().GetRoleSlug(),
		}
	}
	return res
}
