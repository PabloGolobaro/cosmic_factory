package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

// PaymentServer реализует gRPC сервис оплаты.
type PaymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
}

// PayOrder обрабатывает оплату заказа.
func (s *PaymentServer) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	paymentMethod := req.GetPaymentMethod()

	if paymentMethod == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "Платежный метод не выбран")
	}

	_, err := uuid.Parse(req.OrderUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", req.OrderUuid)
	}

	transactionId := uuid.NewString()

	slog.Info("оплата прошла успешно",
		"order_uuid", req.GetOrderUuid(),
		"transaction_id", transactionId,
	)

	return &paymentv1.PayOrderResponse{TransactionUuid: transactionId}, nil
}
