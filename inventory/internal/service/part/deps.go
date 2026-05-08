package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
)

type PartRepository interface {
	Get(ctx context.Context, id uuid.UUID) (entity.Part, error)
	GetBatch(ctx context.Context, ids []uuid.UUID) ([]entity.Part, error)
	GetAll(ctx context.Context) ([]entity.Part, error)
}
