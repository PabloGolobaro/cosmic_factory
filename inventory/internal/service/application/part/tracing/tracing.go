package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	apipart "github.com/PabloGolobaro/cosmic_factory/inventory/internal/api/part/v1"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

const tracerName = "inventory"

type tracedService struct {
	inner apipart.PartService
}

func NewTracedService(inner apipart.PartService) apipart.PartService {
	return &tracedService{inner: inner}
}

func (t *tracedService) Get(ctx context.Context, uuid string) (entity.Part, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "inventory.get_part",
		trace.WithAttributes(
			attribute.String("part.uuid", uuid),
		),
	)
	defer span.End()

	result, err := t.inner.Get(ctx, uuid)
	if err != nil {
		span.RecordError(err)
		return entity.Part{}, err
	}
	return result, nil
}

func (t *tracedService) List(ctx context.Context, uuids []string, partType valueobject.PartType) ([]entity.Part, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "inventory.list_parts",
		trace.WithAttributes(
			attribute.Int("parts.count", len(uuids)),
			attribute.String("parts.type", string(partType)),
		),
	)
	defer span.End()

	result, err := t.inner.List(ctx, uuids, partType)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return result, nil
}

func (t *tracedService) ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "inventory.validate_compatibility",
		trace.WithAttributes(
			attribute.String("slots.hull_uuid", slots.HullUUID),
			attribute.String("slots.engine_uuid", slots.EngineUUID),
		),
	)
	defer span.End()

	if err := t.inner.ValidateCompatibility(ctx, slots); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (t *tracedService) ReserveParts(ctx context.Context, uuids []string) error {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "inventory.reserve_parts",
		trace.WithAttributes(
			attribute.Int("parts.count", len(uuids)),
		),
	)
	defer span.End()

	if err := t.inner.ReserveParts(ctx, uuids); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (t *tracedService) ReleaseParts(ctx context.Context, uuids []string) error {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "inventory.release_parts",
		trace.WithAttributes(
			attribute.Int("parts.count", len(uuids)),
		),
	)
	defer span.End()

	if err := t.inner.ReleaseParts(ctx, uuids); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (t *tracedService) CommitParts(ctx context.Context, uuids []string) error {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "inventory.commit_parts",
		trace.WithAttributes(
			attribute.Int("parts.count", len(uuids)),
		),
	)
	defer span.End()

	if err := t.inner.CommitParts(ctx, uuids); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}
