package main

import (
	"context"
	"log"
	pb "user-grpc/pkg/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

func listUsers(conn *grpc.ClientConn, pageToken string, pageSize int32, sessionId string) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.ListUsersRequest{PageToken: pageToken, PageSize: pageSize}
	response, err := client.ListUsers(attachSessionId(context.Background(), sessionId), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response.GetUsers())
	log.Println(response.GetNextPageToken())
}

func getUser(conn *grpc.ClientConn, id string, sessionId string) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.GetUserRequest{Id: id}
	response, err := client.GetUser(attachSessionId(context.Background(), sessionId), request)

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

func updateUser(conn *grpc.ClientConn, id string, name string, sessionId string) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.UpdateUserRequest{User: &pb.User{Id: id, Email: name}}
	response, err := client.UpdateUser(attachSessionId(context.Background(), sessionId), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response)
}

func deleteUser(conn *grpc.ClientConn, id string, sessionId string) {
	client := pb.NewUserServiceClient(conn)
	request := &pb.DeleteUserRequest{Id: id}
	response, err := client.DeleteUser(attachSessionId(context.Background(), sessionId), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response)
}

func getSesseionId(conn *grpc.ClientConn, email string, password string) string {
	client := pb.NewUserServiceClient(conn)
	request := &pb.SessionAuthUserRequest{User: &pb.User{Email: email, Password: password}}
	response, err := client.AuthUser(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	return response.GetSession()
}

func attachSessionId(ctx context.Context, sessionId string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", sessionId)
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

	//sessionId := getSesseionId(conn, "Tony@maiul.com", "0003")
	//getUser(conn, "10")
	// for i := 0; i < 50; i++ {
	// 	time.Sleep(time.Second * 1)
	// 	listUsers(conn, "1", 5, sessionId)
	// }
	createUser(conn, "Richi@maiul.com", "0004")
	//updateUser(conn, "1", "John")
	//deleteUser(conn, "1")
	//log.Println(sessionId)
}
