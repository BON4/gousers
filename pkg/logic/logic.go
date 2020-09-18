package logic

import (
	"context"
	"time"

	pb "user-grpc/pkg/api"
	rp "user-grpc/repo"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	listTimeOut   = 6
	getTimeOut    = 2
	createTimeOut = 3
	deleteTimeOut = 3
	updateTimeOut = 4
)

type UserServer struct {
	repo        rp.Repository
	sessionRepo rp.SessionRepository
}

//NewUserServer - create new UserServer instace with given repo
func NewUserServer(repo rp.Repository, sessionRepo rp.SessionRepository) pb.UserServiceServer {
	return &UserServer{repo: repo, sessionRepo: sessionRepo}
}

func (u *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*listTimeOut)
	defer cancel()
	resp, err := u.repo.ListUsers(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (u *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if len(req.GetId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid credantials")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*getTimeOut)
	defer cancel()

	resp, err := u.repo.GetUser(ctx, req)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return resp, nil
}

func (u *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	if len(req.GetUser().GetEmail()) == 0 || len(req.GetUser().GetPassword()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid credantials")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*createTimeOut)
	defer cancel()

	resp, err := u.repo.CreateUser(ctx, req)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return resp, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if len(req.GetUser().GetEmail()) == 0 || len(req.GetUser().GetId()) == 0 || len(req.GetUser().GetPassword()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid credantials")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*updateTimeOut)
	defer cancel()

	resp, err := u.repo.UpdateUser(ctx, req)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return resp, nil
}

func (u *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {
	if len(req.GetId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid credantials")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*deleteTimeOut)
	defer cancel()

	resp, err := u.repo.DeleteUser(ctx, req)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return resp, nil
}

func (u *UserServer) AuthUser(ctx context.Context, req *pb.SessionAuthUserRequest) (*pb.SessionAuthUserResponse, error) {
	if len(req.GetUser().GetEmail()) == 0 || len(req.GetUser().GetPassword()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid credantials")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*6)
	defer cancel()

	ifUser, err := u.repo.PasswordCheck(ctx, req)

	if err != nil || ifUser == false {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	resp, err := u.sessionRepo.CreateSession(ctx, req)

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, nil
}
