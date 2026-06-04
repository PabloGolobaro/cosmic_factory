package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	svcauth "github.com/PabloGolobaro/cosmic_factory/iam/internal/service/auth"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
)

const tracerName = "iam.auth"

type tracedService struct {
	inner svcauth.Service
}

func NewTracedService(inner svcauth.Service) svcauth.Service {
	return &tracedService{inner: inner}
}

func (t *tracedService) Login(ctx context.Context, in input.LoginInput) (model.Session, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "iam.auth.login",
		trace.WithAttributes(
			attribute.String("user.login", in.Login),
		),
	)
	defer span.End()

	session, err := t.inner.Login(ctx, in)
	if err != nil {
		span.RecordError(err)
		return model.Session{}, err
	}
	return session, nil
}

func (t *tracedService) Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "iam.auth.whoami",
		trace.WithAttributes(
			attribute.String("session.uuid", sessionUUID),
		),
	)
	defer span.End()

	session, user, err := t.inner.Whoami(ctx, sessionUUID)
	if err != nil {
		span.RecordError(err)
		return model.Session{}, model.User{}, err
	}
	return session, user, nil
}

func (t *tracedService) Logout(ctx context.Context, sessionUUID string) error {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "iam.auth.logout",
		trace.WithAttributes(
			attribute.String("session.uuid", sessionUUID),
		),
	)
	defer span.End()

	if err := t.inner.Logout(ctx, sessionUUID); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}
