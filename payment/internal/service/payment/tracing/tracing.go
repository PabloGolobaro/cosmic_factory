package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	apipayment "github.com/PabloGolobaro/cosmic_factory/payment/internal/api/payment/v1"
)

const tracerName = "payment"

type tracedService struct {
	inner apipayment.PaymentService
}

func NewTracedService(inner apipayment.PaymentService) apipayment.PaymentService {
	return &tracedService{inner: inner}
}

func (t *tracedService) Pay(ctx context.Context, uuid, paymentMethod string) (string, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "payment.pay",
		trace.WithAttributes(
			attribute.String("payment.order_uuid", uuid),
			attribute.String("payment.method", paymentMethod),
		),
	)
	defer span.End()

	txUUID, err := t.inner.Pay(ctx, uuid, paymentMethod)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	span.SetAttributes(attribute.String("payment.transaction_uuid", txUUID))
	return txUUID, nil
}
