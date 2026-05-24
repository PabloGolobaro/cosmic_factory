package session

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/converter"
)

func (s *store) Save(ctx context.Context, session model.Session) error {
	key := sessionKey(session.UUID)
	view := converter.SessionToRedisView(session)

	_, err := s.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, view)
		pipe.Expire(ctx, key, s.ttl)
		return nil
	})
	if err != nil {
		return fmt.Errorf("сохранить сессию: %w", err)
	}

	return nil
}
