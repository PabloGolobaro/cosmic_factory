package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

func (s *store) GetByUUID(ctx context.Context, id uuid.UUID) (model.User, error) {
	return s.queryUser(ctx, selectUserByUUID, id)
}
