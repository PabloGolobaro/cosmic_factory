package app

import (
	"log/slog"

	"buf.build/go/protovalidate"
	protovalidateMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	v1 "github.com/PabloGolobaro/cosmic_factory/inventory/internal/api/part/v1"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/part"
	partSvc "github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/application/part"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/domain"
	"github.com/PabloGolobaro/cosmic_factory/shared/pkg/interceptors"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool) {
	store := part.NewPartStore(pool)
	svc := partSvc.NewPartService(store, domain.NewCompatibilityChecker())
	api := v1.NewApi(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}

func Interceptors() []grpc.ServerOption {
	validator, err := protovalidate.New()
	if err != nil {
		slog.Error("ошибка создания валидатора", "error", err)
	}
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryInterceptor(),
			interceptors.LoggerInterceptor(),
			protovalidateMiddleware.UnaryServerInterceptor(validator),
		),
	}
}
