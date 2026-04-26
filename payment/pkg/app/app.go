package app

import (
	"log/slog"

	"buf.build/go/protovalidate"
	protovalidateMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"

	v1 "github.com/PabloGolobaro/cosmic_factory/payment/internal/api/payment/v1"
	"github.com/PabloGolobaro/cosmic_factory/payment/internal/service/payment"
	"github.com/PabloGolobaro/cosmic_factory/shared/pkg/interceptors"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

func RegisterServices(grpcServer *grpc.Server) {
	svc := payment.NewPaymentService()
	api := v1.NewApi(svc)
	paymentv1.RegisterPaymentServiceServer(grpcServer, api)
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
