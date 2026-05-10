package v1

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

type PartService interface {
	Get(context.Context, string) (entity.Part, error)
	List(context.Context, []string, valueobject.PartType) ([]entity.Part, error)
	ValidateCompatibility(ctx context.Context, slots valueobject.ShipSlots) error
	ReserveParts(ctx context.Context, uuids []string) error
	ReleaseParts(ctx context.Context, uuids []string) error
}
