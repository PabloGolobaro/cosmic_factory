package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/converter"
)

const (
	insertUser           = `INSERT INTO users (uuid, login, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	pgErrUniqueViolation = "23505"
)

func (s *store) Create(ctx context.Context, user model.User) error {
	rec := converter.UserToRecord(user)

	_, err := s.pool.Exec(ctx, insertUser, rec.UUID, rec.Login, rec.PasswordHash, rec.CreatedAt, rec.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
			return errs.ErrUserAlreadyExists
		}

		return fmt.Errorf("создать пользователя: %w", err)
	}

	return nil
}
