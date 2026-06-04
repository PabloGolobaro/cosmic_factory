package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	orderapi "github.com/PabloGolobaro/cosmic_factory/order/internal/api/order/v1"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

const tracerName = "order"

type tracedService struct {
	inner orderapi.OrderService
}

func NewTracedService(inner orderapi.OrderService) orderapi.OrderService {
	return &tracedService{inner: inner}
}

func (t *tracedService) Create(ctx context.Context, order model.Order) (model.Order, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "order.create",
		trace.WithAttributes(
			attribute.String("order.hull_uuid", order.HullUUID.String()),
			attribute.String("order.engine_uuid", order.EngineUUID.String()),
		),
	)
	defer span.End()

	result, err := t.inner.Create(ctx, order)
	if err != nil {
		span.RecordError(err)
		return model.Order{}, err
	}

	span.SetAttributes(
		attribute.String("order.uuid", result.OrderUUID.String()),
		attribute.Int64("order.total_price", result.TotalPrice),
	)
	return result, nil
}

func (t *tracedService) Get(ctx context.Context, id string) (*model.Order, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "order.get",
		trace.WithAttributes(
			attribute.String("order.uuid", id),
		),
	)
	defer span.End()

	result, err := t.inner.Get(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return result, nil
}

func (t *tracedService) Cancel(ctx context.Context, id string) error {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "order.cancel",
		trace.WithAttributes(
			attribute.String("order.uuid", id),
		),
	)
	defer span.End()

	if err := t.inner.Cancel(ctx, id); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (t *tracedService) Pay(ctx context.Context, id string, method model.PaymentMethod) (string, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "order.pay",
		trace.WithAttributes(
			attribute.String("order.uuid", id),
			attribute.Int64("order.payment_method", int64(method)),
		),
	)
	defer span.End()

	txUUID, err := t.inner.Pay(ctx, id, method)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	span.SetAttributes(attribute.String("order.transaction_uuid", txUUID))
	return txUUID, nil
}
