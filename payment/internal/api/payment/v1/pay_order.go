package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PabloGolobaro/cosmic_factory/payment/internal/converter"
	errs "github.com/PabloGolobaro/cosmic_factory/payment/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/payment/internal/model"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	paymentMethod := converter.PaymentMethodToString(model.PaymentMethod(req.GetPaymentMethod()))

	if err := a.PaymentService.Pay(ctx, req.GetOrderUuid(), paymentMethod); err != nil {
		if errors.Is(err, errs.ErrInvalidUUID) || errors.Is(err, errs.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &paymentv1.PayOrderResponse{}, nil
}
