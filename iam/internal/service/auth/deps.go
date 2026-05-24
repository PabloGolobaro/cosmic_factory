package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

type userRepository interface {
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetByUUID(ctx context.Context, id uuid.UUID) (model.User, error)
}

type sessionRepository interface {
	Save(ctx context.Context, session model.Session) error
	Get(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error)
	Delete(ctx context.Context, sessionUUID uuid.UUID) error
}
