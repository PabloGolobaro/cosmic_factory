package order

import (
	"sync"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/record"
)

type store struct {
	mu     sync.RWMutex
	orders map[string]record.OrderRecord
}

func NewOrderStore() *store {
	return &store{
		orders: make(map[string]record.OrderRecord),
	}
}
