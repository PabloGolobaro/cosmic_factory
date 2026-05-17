package order

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	pool *pgxpool.Pool

	getter *trmpgx.CtxGetter
}

func NewOrderRepo(pool *pgxpool.Pool) *repo {
	return &repo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

// New — алиас NewOrderRepo для использования в e2e-тестах.
// txManager игнорируется: DefaultCtxGetter работает с любым manager.Manager.
func New(pool *pgxpool.Pool, _ *manager.Manager) *repo {
	return NewOrderRepo(pool)
}
