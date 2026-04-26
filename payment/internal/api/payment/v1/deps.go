package v1

import "context"

type PaymentService interface {
	Pay(ctx context.Context, uuid, paymentMethod string) (string, error)
}
