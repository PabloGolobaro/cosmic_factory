package orderitem

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
)

func (s *repo) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM order_items WHERE uuid = $1`

	cmdTag, err := s.getter.DefaultTrOrDB(ctx, s.pool).Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %s", errs.ErrOrderNotFound, id)
	}
	return nil
}
