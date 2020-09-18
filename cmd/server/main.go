package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	pb "user-grpc/pkg/api"
	"user-grpc/pkg/logic"
	rp "user-grpc/repo"

	"github.com/go-redis/redis/v8"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

var (
	logrusLogger *logrus.Logger
	customFunc   grpc_logrus.CodeToLevel
)

func main() {
	addr := "0.0.0.0:8080"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
		//panic(err)
	}

	// logrusEntry := logrus.NewEntry(logrusLogger)

	// opts := []grpc_logrus.Option{
	// 	grpc_logrus.WithLevels(customFunc),
	// }

	// grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	// s := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
	// 	grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
	// 	grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
	// ),
	// )

	s := grpc.NewServer()

	// TODO: Replace with your own certificate!
	//grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),

	//Move sql open somewhere
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		log.Fatal(err)
		//panic(err)
	}

	repo := rp.NewPsqlREpository(db)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		log.Fatal(status.Err())
	}

	defer rdb.Close()

	redisRepo := rp.NewRedisRepository(rdb)

	pb.RegisterUserServiceServer(s, logic.NewUserServer(repo, redisRepo))

	// Serve gRPC Server
	log.Println("Serving gRPC on https://", addr)
	log.Fatal(s.Serve(lis))
	// if err := s.Serve(lis); err != nil {
	// 	panic(err)
	// }
}
