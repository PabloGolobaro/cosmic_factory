package v1

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
)

type PartService interface {
	Get(context.Context, string) (model.Part, error)
	List(context.Context, []string, model.PartType) ([]model.Part, error)
}
