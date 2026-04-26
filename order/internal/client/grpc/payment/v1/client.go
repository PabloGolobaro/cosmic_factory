package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

type paymentClient struct {
	client paymentv1.PaymentServiceClient
}

func NewPaymentClient(client paymentv1.PaymentServiceClient) *paymentClient {
	return &paymentClient{client: client}
}

func (p paymentClient) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) error {
	_, err := p.client.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: paymentMethodToProto(paymentMethod),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return errs.ErrOrderNotFound
			case codes.FailedPrecondition:
				return errs.ErrOrderAlreadyPaid
			case codes.InvalidArgument:
				return errs.ErrInvalidPaymentMethod
			}
		}
		return fmt.Errorf("оплатить заказ: %w", err)
	}
	return nil
}
