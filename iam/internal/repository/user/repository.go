package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/record"
)

const (
	selectUserCols    = `SELECT uuid, login, password_hash, created_at, updated_at FROM users`
	selectUserByLogin = selectUserCols + ` WHERE login = $1`
	selectUserByUUID  = selectUserCols + ` WHERE uuid = $1`
)

type Repository interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetByUUID(ctx context.Context, id uuid.UUID) (model.User, error)
}

type store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repository {
	return &store{pool: pool}
}

func (s *store) queryUser(ctx context.Context, sql string, arg any) (model.User, error) {
	rec := record.UserRecord{}
	err := s.pool.QueryRow(ctx, sql, arg).Scan(
		&rec.UUID, &rec.Login, &rec.PasswordHash, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errs.ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("запрос пользователя: %w", err)
	}

	u, err := converter.UserFromRecord(rec)
	if err != nil {
		return model.User{}, fmt.Errorf("конвертировать запись: %w", err)
	}

	return u, nil
}
