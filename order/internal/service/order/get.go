package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s service) Get(ctx context.Context, id string) (*model.Order, error) {
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}
	order, err := s.Repository.Get(ctx, orderUUID)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
