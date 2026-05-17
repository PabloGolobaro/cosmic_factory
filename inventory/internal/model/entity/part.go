package entity

import (
	"fmt"
	"time"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

// Part — доменная сущность детали космического корабля.
type Part struct {
	uuid          string
	name          string
	description   string
	partType      valueobject.PartType
	price         int64
	stockQuantity int
	reserved      int
	properties    valueobject.PartProperties
	createdAt     time.Time
}

func (p *Part) UUID() string                           { return p.uuid }
func (p *Part) Name() string                           { return p.name }
func (p *Part) Description() string                    { return p.description }
func (p *Part) PartType() valueobject.PartType         { return p.partType }
func (p *Part) Price() int64                           { return p.price }
func (p *Part) StockQuantity() int64                   { return int64(p.stockQuantity) }
func (p *Part) Reserved() int                          { return p.reserved }
func (p *Part) Properties() valueobject.PartProperties { return p.properties }
func (p *Part) CreatedAt() time.Time                   { return p.createdAt }

// Reserve резервирует n единиц. Ошибка если reserved+n > stockQuantity.
func (p *Part) Reserve(n int) error {
	if p.reserved+n > p.stockQuantity {
		return fmt.Errorf("%w: доступно %d, запрошено %d", errs.ErrOutOfStock, p.stockQuantity-p.reserved, n)
	}
	p.reserved += n
	return nil
}

// Release освобождает n единиц резерва. Ошибка если reserved-n < 0.
func (p *Part) Release(n int) error {
	if p.reserved-n < 0 {
		return fmt.Errorf("%w: зарезервировано %d, попытка освободить %d", errs.ErrNothingToRelease, p.reserved, n)
	}
	p.reserved -= n
	return nil
}

// Commit списывает 1 единицу со склада. Ошибка если stock или reserved равны 0.
func (p *Part) Commit() error {
	if p.stockQuantity <= 0 {
		return fmt.Errorf("%w: stock равен 0", errs.ErrNothingToCommit)
	}
	if p.reserved <= 0 {
		return fmt.Errorf("%w: reserved равен 0", errs.ErrNothingToCommit)
	}
	p.stockQuantity--
	p.reserved--
	return nil
}

// RestorePart восстанавливает сущность из БД (без валидации — данные уже проверены).
func RestorePart(uuid, name, description string, partType valueobject.PartType, price int64,
	stockQuantity, reserved int, properties valueobject.PartProperties, createdAt time.Time,
) Part {
	return Part{
		uuid:          uuid,
		name:          name,
		description:   description,
		partType:      partType,
		price:         price,
		stockQuantity: stockQuantity,
		reserved:      reserved,
		properties:    properties,
		createdAt:     createdAt,
	}
}
