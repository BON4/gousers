package repo

import (
	"context"
	"fmt"
	"time"
	pb "user-grpc/pkg/api"

	redis "github.com/go-redis/redis/v8"
	uuid "github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RedisRepository struct {
	conn *redis.Client
}

func NewRedisRepository(cache *redis.Client) SessionRepository {
	return &RedisRepository{conn: cache}
}

func (r *RedisRepository) CreateSession(ctx context.Context, req *pb.SessionAuthUserRequest) (*pb.SessionAuthUserResponse, error) {
	var userKey string = fmt.Sprintf("%s:%s", req.GetUser().GetId(), req.GetUser().GetEmail())
	var userSessionId string = uuid.New().String()
	redisResp := r.conn.Set(ctx, userSessionId, userKey, time.Second*10)
	if redisResp.Err() != nil {
		return nil, redisResp.Err()
	}

	sessionRep := pb.SessionAuthUserResponse{Session: userSessionId}
	return &sessionRep, nil
}

func (r *RedisRepository) VerifySession(ctx context.Context, sessionId string) (bool, error) {
	redisResp := r.conn.Get(ctx, sessionId)

	if redisResp.Err() != nil {
		return false, status.Error(codes.Unauthenticated, "invalid session Id")
	}

	if redisResp == nil {
		return false, nil
	}

	return true, nil
}
