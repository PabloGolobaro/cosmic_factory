package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

func (s *service) GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	return s.userRepo.GetByUUID(ctx, userUUID)
}
