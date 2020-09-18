package main

import (
	"context"
	"log"
	pb "user-grpc/pkg/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func listUsers(conn *grpc.ClientConn, pageToken string, pageSize int32) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.ListUsersRequest{PageToken: pageToken, PageSize: pageSize}
	response, err := client.ListUsers(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response.GetUsers())
	log.Println(response.GetNextPageToken())
}

func getUser(conn *grpc.ClientConn, id string) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.GetUserRequest{Id: id}
	response, err := client.GetUser(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response)
}

func createUser(conn *grpc.ClientConn, email string, password string) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.CreateUserRequest{User: &pb.User{Email: email, Password: password}}
	response, err := client.CreateUser(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response)
}

func updateUser(conn *grpc.ClientConn, id string, name string) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.UpdateUserRequest{User: &pb.User{Id: id, Email: name}}
	response, err := client.UpdateUser(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response)
}

func deleteUser(conn *grpc.ClientConn, id string) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.DeleteUserRequest{Id: id}
	response, err := client.DeleteUser(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response)
}

func getSesseionId(conn *grpc.ClientConn, email string, password string) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.SessionAuthUserRequest{User: &pb.User{Email: email, Password: password}}
	response, err := client.AuthUser(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response)
}

func reduceLeftInt(arr []int, first int, f func(total int, rest ...int) int) int {
	for _, r := range arr {
		first = f(first, r)
	}
	return first
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("0.0.0.0:8080", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//getUser(conn, "10")
	//listUsers(conn, "1", 5)
	//createUser(conn, "Jake@mail.com", "1488")
	//updateUser(conn, "1", "John")
	//deleteUser(conn, "1")
	getSesseionId(conn, "Jake@mail.com", "1488")
}
