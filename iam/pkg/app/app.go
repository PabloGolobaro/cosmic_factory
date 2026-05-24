package app

import (
	"fmt"
	"time"

	"buf.build/go/protovalidate"
	protovalidateMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiauth "github.com/PabloGolobaro/cosmic_factory/iam/internal/api/auth/v1"
	apiuser "github.com/PabloGolobaro/cosmic_factory/iam/internal/api/user/v1"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/interceptor"
	reposes "github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/session"
	repouser "github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/user"
	svcauth "github.com/PabloGolobaro/cosmic_factory/iam/internal/service/auth"
	svcuser "github.com/PabloGolobaro/cosmic_factory/iam/internal/service/user"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/grpc/health"
	"github.com/PabloGolobaro/cosmic_factory/shared/pkg/interceptors"
	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
	userproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/user/v1"
)

// NewGRPCServer создаёт gRPC-сервер IAM с полным стеком зависимостей.
// Используется в интеграционных тестах — testcontainers предоставляет pool и redisClient.
func NewGRPCServer(pool *pgxpool.Pool, redisClient *redis.Client, sessionTTL time.Duration, bcryptCost int) (*grpc.Server, error) {
	userRepo := repouser.New(pool)
	sessionRepo := reposes.New(redisClient, sessionTTL)

	userSvc := svcuser.New(userRepo, bcryptCost)
	authSvc := svcauth.New(userRepo, sessionRepo, sessionTTL)

	authHandler := apiauth.NewAPI(authSvc)
	userHandler := apiuser.NewAPI(userSvc)

	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("создание protovalidate валидатора: %w", err)
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryInterceptor(),
			interceptors.LoggerInterceptor(),
			protovalidateMiddleware.UnaryServerInterceptor(validator),
			interceptor.ErrorInterceptor(),
		),
	)

	authproto.RegisterAuthServiceServer(srv, authHandler)
	userproto.RegisterUserServiceServer(srv, userHandler)
	health.RegisterService(srv)
	reflection.Register(srv)

	return srv, nil
}
