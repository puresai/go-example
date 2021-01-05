package endpoint

import (
	"context"
	"demo7-docker-compose/service"

	"github.com/go-kit/kit/endpoint"
)

type UserEndpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginRes struct {
	UserInfo *service.UserInfoDTO `json:"userinfo"`
}

func MakeLoginEndpoint(userS service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (res interface{}, err error) {
		req := request.(*LoginRequest)
		userinfo, err := userS.Login(ctx, req.Email, req.Password)
		return &LoginRes{userinfo}, err
	}
}

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

type RegisterRes struct {
	UserInfo *service.UserInfoDTO `json:"userinfo"`
}

func MakeRegisterEndpoint(userS service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*RegisterRequest)
		userInfo, err := userS.Register(ctx, &service.RegisterUser{
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
		})
		return &RegisterRes{UserInfo: userInfo}, err

	}
}
