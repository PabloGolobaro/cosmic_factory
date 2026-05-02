package orderitem

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewOrderItemRepo(pool *pgxpool.Pool) *repo {
	return &repo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}
