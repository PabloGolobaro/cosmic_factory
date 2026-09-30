package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
)

func (s *service) Authorize(ctx context.Context, in input.AuthorizeInput) (model.AuthzDecision, error) {
	_, user, err := s.Whoami(ctx, in.SessionUUID)
	if err != nil {
		return model.AuthzDecision{}, err
	}

	allowed, err := s.policy.Allow(ctx, model.AuthzInput{
		Subject:  model.AuthzSubject{UUID: user.UUID.String(), Role: user.Role},
		Action:   in.Action,
		Resource: model.AuthzResource{OwnerUUID: in.OwnerUUID},
	})
	if err != nil {
		return model.AuthzDecision{}, fmt.Errorf("проверка политики: %w", err)
	}

	// Decision log: every authorization decision is recorded for audit.
	slog.InfoContext(ctx, "решение авторизации",
		"user_uuid", user.UUID.String(),
		"role", user.Role,
		"action", in.Action,
		"owner_uuid", in.OwnerUUID,
		"allowed", allowed,
	)

	return model.AuthzDecision{Allowed: allowed, User: user}, nil
}
