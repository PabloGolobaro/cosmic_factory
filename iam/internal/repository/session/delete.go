package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *store) Delete(ctx context.Context, sessionUUID uuid.UUID) error {
	if err := s.client.Del(ctx, sessionKey(sessionUUID)).Err(); err != nil {
		return fmt.Errorf("удалить сессию: %w", err)
	}

	return nil
}
