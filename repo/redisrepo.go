package repo

import (
	"context"
	"fmt"
	"time"
	pb "user-grpc/pkg/api"

	redis "github.com/go-redis/redis/v8"
	uuid "github.com/google/uuid"
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
	redisResp := r.conn.Set(ctx, userKey, userSessionId, time.Second*30)
	if redisResp.Err() != nil {
		return nil, redisResp.Err()
	}

	sessionRep := pb.SessionAuthUserResponse{Session: userSessionId}
	return &sessionRep, nil
}
