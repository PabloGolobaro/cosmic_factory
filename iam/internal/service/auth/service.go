package auth

import (
	"context"
	"time"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
)

type Service interface {
	Login(ctx context.Context, in input.LoginInput) (model.Session, error)
	Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error)
	Logout(ctx context.Context, sessionUUID string) error
	Authorize(ctx context.Context, in input.AuthorizeInput) (model.AuthzDecision, error)
}

type service struct {
	userRepo    userRepository
	sessionRepo sessionRepository
	policy      policyEngine
	ttl         time.Duration
}

func New(userRepo userRepository, sessionRepo sessionRepository, policy policyEngine, ttl time.Duration) Service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		policy:      policy,
		ttl:         ttl,
	}
}
