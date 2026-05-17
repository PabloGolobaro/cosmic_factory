package assembly

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/model"
)

type Producer interface {
	PublishShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
