package auth

import (
	"context"
	"microservices/auth/internal/modules/auth/service"

	pb "microservices/auth/authgrpc/auth"
)

type AuthServiceRPC struct {
	authService service.Auther
}

func NewAuthServiceRPC(authService service.Auther) *AuthServiceRPC {
	return &AuthServiceRPC{authService: authService}
}

func (a *AuthServiceRPC) Register(in service.RegisterIn, out *service.RegisterOut) error {
	*out = a.authService.Register(in)

	return nil
}

func (a *AuthServiceRPC) Login(in service.LoginIn, out *service.LoginOut) error {
	*out = a.authService.Login(in)

	return nil
}

type AuthServiceGRPC struct {
	authService service.Auther
	pb.UnimplementedAutherServer
}

func NewUserServiceGRPC(userService service.Auther) *AuthServiceGRPC {
	return &AuthServiceGRPC{
		authService: userService,
	}
}

func (a *AuthServiceGRPC) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	res := a.authService.Login(service.LoginIn{Ctx: ctx, Email: in.Email, Password: in.Password})

	return &pb.LoginResponse{Success: res.Success, Message: res.Message}, nil
}

func (a *AuthServiceGRPC) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	res := a.authService.Register(service.RegisterIn{Name: in.Name, Email: in.Email, Password: in.Password, Phone: in.Phone})
	if res.Error != nil {
		return &pb.RegisterResponse{
			Status:  uint32(res.Status),
			Message: res.Error.Error(),
		}, res.Error
	}
	return &pb.RegisterResponse{Status: uint32(res.Status)}, nil
}
