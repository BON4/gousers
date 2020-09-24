package logic

import (
	"context"
	"log"
	rp "user-grpc/repo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthUserInterceptor struct {
	sessionRepo      rp.SessionRepository
	accesibleMethods []string
}

func NewAuthUserInterceptor(redisSession rp.SessionRepository, accesibleMethods []string) *AuthUserInterceptor {
	return &AuthUserInterceptor{sessionRepo: redisSession, accesibleMethods: accesibleMethods}
}

func (interceptor *AuthUserInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthUserInterceptor) authorize(ctx context.Context, method string) error {
	if !contains(interceptor.accesibleMethods, method) {
		//Its a public method
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]

	if len(values) == 0 {
		return status.Error(codes.Unauthenticated, "session ID is not provided")
	}

	sessionId := values[0]
	ok, err := interceptor.sessionRepo.VerifySession(ctx, sessionId)
	if err != nil {
		return err
	}

	if !ok {
		return status.Error(codes.Unauthenticated, "invalid credetials")
	} else {
		return nil
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
