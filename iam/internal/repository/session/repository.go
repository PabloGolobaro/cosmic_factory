package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

type Repository interface {
	Save(ctx context.Context, session model.Session) error
	Get(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error)
	Delete(ctx context.Context, sessionUUID uuid.UUID) error
}

type store struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) Repository {
	return &store{client: client, ttl: ttl}
}

func sessionKey(id uuid.UUID) string {
	return "session:" + id.String()
}
