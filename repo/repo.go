package repo

import (
	"context"

	pb "user-grpc/pkg/api"

	"github.com/golang/protobuf/ptypes/empty"
)

type Repository interface {
	ListUsers(context.Context, *pb.ListUsersRequest) (*pb.ListUsersResponse, error)
	GetUser(context.Context, *pb.GetUserRequest) (*pb.User, error)
	CreateUser(context.Context, *pb.CreateUserRequest) (*pb.User, error)
	UpdateUser(context.Context, *pb.UpdateUserRequest) (*pb.User, error)
	DeleteUser(context.Context, *pb.DeleteUserRequest) (*empty.Empty, error)
	PasswordCheck(context.Context, *pb.SessionAuthUserRequest) (bool, error)
}

type SessionRepository interface {
	CreateSession(context.Context, *pb.SessionAuthUserRequest) (*pb.SessionAuthUserResponse, error)
	//TODO
	//UpdateSession(context.Context, ...) (...)
}
