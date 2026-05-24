package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

type userRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByUUID(ctx context.Context, id uuid.UUID) (model.User, error)
}
