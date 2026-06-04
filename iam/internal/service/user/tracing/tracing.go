package tracing

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
	svcuser "github.com/PabloGolobaro/cosmic_factory/iam/internal/service/user"
)

const tracerName = "iam.user"

type tracedService struct {
	inner svcuser.Service
}

func NewTracedService(inner svcuser.Service) svcuser.Service {
	return &tracedService{inner: inner}
}

func (t *tracedService) Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "iam.user.register",
		trace.WithAttributes(
			attribute.String("user.login", in.Login),
		),
	)
	defer span.End()

	userUUID, err := t.inner.Register(ctx, in)
	if err != nil {
		span.RecordError(err)
		return uuid.UUID{}, err
	}

	span.SetAttributes(attribute.String("user.uuid", userUUID.String()))
	return userUUID, nil
}

func (t *tracedService) GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "iam.user.get",
		trace.WithAttributes(
			attribute.String("user.uuid", userUUID.String()),
		),
	)
	defer span.End()

	user, err := t.inner.GetUser(ctx, userUUID)
	if err != nil {
		span.RecordError(err)
		return model.User{}, err
	}
	return user, nil
}
