package user

import (
	"context"
	"microservices/user/internal/models"
	"microservices/user/internal/modules/user/service"

	pb "microservices/user/usergrpc/user"
)

type UserServiceRPC struct {
	userService service.Userer
}

func NewUserService(userService service.Userer) *UserServiceRPC {
	return &UserServiceRPC{userService: userService}
}

func (u *UserServiceRPC) Profile(in ProfileIn, out *ProfileOut) error {

	user, err := u.userService.Profile(in.Email)
	if err != nil {
		return err
	}

	*out = ProfileOut{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	return nil
}

func (u *UserServiceRPC) List(in ListIn, out *ListOut) error {
	users, err := u.userService.List()
	if err != nil {
		return err
	}

	*out = ListOut{
		Users: users,
	}

	return nil
}

func (u *UserServiceRPC) Create(in CreateIn, out *CreateOut) error {
	err := u.userService.Create(models.UserDTO{Name: in.Name, Email: in.Email, Password: in.Password})
	if err != nil {
		return err
	}

	*out = CreateOut{
		Success: true,
	}

	return nil
}

type UserServiceGRPC struct {
	userService service.Userer
	pb.UnimplementedUsererServer
}

func NewUserServiceGRPC(userService service.Userer) *UserServiceGRPC {
	return &UserServiceGRPC{
		userService: userService,
	}
}

func (g *UserServiceGRPC) Profile(ctx context.Context, in *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	user, err := g.userService.Profile(in.Email)
	if err != nil {
		return nil, err
	}

	return &pb.ProfileResponse{Id: uint32(user.ID), Name: user.Name, Email: user.Email, Passwrd: user.Password, Phone: user.Phone}, nil
}

func (g *UserServiceGRPC) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	list, err := g.userService.List()
	if err != nil {
		return nil, err
	}

	users := make([]*pb.User, len(list))
	for i, user := range list {
		users[i] = &pb.User{
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
		}
	}

	return &pb.ListResponse{Users: users}, nil
}

func (g *UserServiceGRPC) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	err := g.userService.Create(models.UserDTO{Name: in.Name, Email: in.Email, Password: in.Password, Phone: in.Phone})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Success: true}, nil
}

func (g *UserServiceGRPC) GetByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	res, err := g.userService.GeyByID(ctx, int(in.GetId()))
	if err != nil {
		return nil, err
	}

	return &pb.GetByIDResponse{Email: res.Email, Phone: res.Phone}, nil
}
