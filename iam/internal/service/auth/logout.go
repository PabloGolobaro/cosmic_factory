package auth

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) Logout(ctx context.Context, sessionUUID string) error {
	id, err := uuid.Parse(sessionUUID)
	if err != nil {
		return nil
	}

	return s.sessionRepo.Delete(ctx, id)
}
