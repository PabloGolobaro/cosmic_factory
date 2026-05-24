package app

import (
	"log/slog"

	"buf.build/go/protovalidate"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	protovalidateMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	v1 "github.com/PabloGolobaro/cosmic_factory/inventory/internal/api/part/v1"
	iamv1 "github.com/PabloGolobaro/cosmic_factory/inventory/internal/client/grpc/iam/v1"
	authinterceptor "github.com/PabloGolobaro/cosmic_factory/inventory/internal/interceptor"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/part"
	partSvc "github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/application/part"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/domain"
	"github.com/PabloGolobaro/cosmic_factory/shared/pkg/interceptors"
	authv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool) {
	txm, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		slog.Error("ошибка создания transaction manager", "error", err)
	}
	store := part.NewPartStore(pool)
	svc := partSvc.NewPartService(store, domain.NewCompatibilityChecker(), txm)
	api := v1.NewApi(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}

func Interceptors(authClients ...authv1.AuthServiceClient) []grpc.ServerOption {
	validator, err := protovalidate.New()
	if err != nil {
		slog.Error("ошибка создания валидатора", "error", err)
	}
	chain := make([]grpc.UnaryServerInterceptor, 0, 4)
	chain = append(chain,
		interceptors.RecoveryInterceptor(),
		interceptors.LoggerInterceptor(),
		protovalidateMiddleware.UnaryServerInterceptor(validator),
	)
	if len(authClients) > 0 {
		chain = append(chain, authinterceptor.New(iamv1.New(authClients[0])))
	}
	return []grpc.ServerOption{grpc.ChainUnaryInterceptor(chain...)}
}
