package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/redis_view"
)

func (s *store) Get(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error) {
	key := sessionKey(sessionUUID)

	var view redis_view.SessionRedisView
	if err := s.client.HGetAll(ctx, key).Scan(&view); err != nil {
		return model.Session{}, fmt.Errorf("получить сессию: %w", err)
	}

	if view.UUID == "" {
		return model.Session{}, errs.ErrSessionNotFound
	}

	session, err := converter.SessionFromRedisView(view)
	if err != nil {
		return model.Session{}, fmt.Errorf("конвертировать сессию: %w", err)
	}

	return session, nil
}
