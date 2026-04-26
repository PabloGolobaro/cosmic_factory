package app

import (
	"log/slog"

	"buf.build/go/protovalidate"
	protovalidateMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"

	v1 "github.com/PabloGolobaro/cosmic_factory/inventory/internal/api/part/v1"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/part"
	partSvc "github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/part"
	"github.com/PabloGolobaro/cosmic_factory/shared/pkg/interceptors"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

func RegisterServices(grpcServer *grpc.Server) {
	store := part.NewPartStore()
	svc := partSvc.NewPartService(store)
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
