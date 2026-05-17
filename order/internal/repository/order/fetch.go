package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/record"
)

func (s *repo) fetchOrder(ctx context.Context, id uuid.UUID, sql string) (model.Order, error) {
	orderRecord := record.OrderRecord{}

	err := s.getter.DefaultTrOrDB(ctx, s.pool).QueryRow(ctx, sql, id).Scan(
		&orderRecord.OrderUUID, &orderRecord.TotalPrice, &orderRecord.Status,
		&orderRecord.TransactionUUID, &orderRecord.PaymentMethod,
		&orderRecord.CreatedAt, &orderRecord.UpdatedAt, &orderRecord.UserUUID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, fmt.Errorf("%w: %s", errs.ErrOrderNotFound, id)
		}

		return model.Order{}, err
	}

	order := converter.OrderFromRecord(orderRecord)

	itemSQL := `SELECT * FROM order_items WHERE order_uuid = $1`

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, itemSQL, id)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		r := record.OrderItemRecord{}
		if err = rows.Scan(&r.UUID, &r.OrderUUID, &r.PartUUID,
			&r.PartType, &r.Price, &r.CreatedAt); err != nil {
			return model.Order{}, err
		}
		item := converter.OrderItemFromRecord(r)
		switch item.PartType {
		case model.PartTypeHull:
			order.HullUUID = item.PartUUID
		case model.PartTypeEngine:
			order.EngineUUID = item.PartUUID
		case model.PartTypeShield:
			partUUID := item.PartUUID
			order.ShieldUUID = &partUUID
		case model.PartTypeWeapon:
			partUUID := item.PartUUID
			order.WeaponUUID = &partUUID
		}
	}
	if err = rows.Err(); err != nil {
		return model.Order{}, err
	}

	return order, nil
}
