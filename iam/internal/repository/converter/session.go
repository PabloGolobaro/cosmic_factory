package converter

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/redis_view"
)

func SessionToRedisView(s model.Session) redis_view.SessionRedisView {
	return redis_view.SessionRedisView{
		UUID:      s.UUID.String(),
		UserUUID:  s.UserUUID.String(),
		Login:     s.Login,
		CreatedAt: s.CreatedAt.UTC().Format(time.RFC3339),
		ExpiresAt: s.ExpiresAt.UTC().Format(time.RFC3339),
	}
}

func SessionFromRedisView(v redis_view.SessionRedisView) (model.Session, error) {
	id, err := uuid.Parse(v.UUID)
	if err != nil {
		return model.Session{}, fmt.Errorf("некорректный UUID сессии: %w", err)
	}

	userID, err := uuid.Parse(v.UserUUID)
	if err != nil {
		return model.Session{}, fmt.Errorf("некорректный UUID пользователя: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, v.CreatedAt)
	if err != nil {
		return model.Session{}, fmt.Errorf("некорректный created_at: %w", err)
	}

	expiresAt, err := time.Parse(time.RFC3339, v.ExpiresAt)
	if err != nil {
		return model.Session{}, fmt.Errorf("некорректный expires_at: %w", err)
	}

	return model.Session{
		UUID:      id,
		UserUUID:  userID,
		Login:     v.Login,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}
