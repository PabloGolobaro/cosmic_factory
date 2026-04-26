package v1

import (
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

type api struct {
	PaymentService PaymentService
	paymentv1.UnimplementedPaymentServiceServer
}

func NewApi(paymentService PaymentService) *api {
	return &api{PaymentService: paymentService}
}
