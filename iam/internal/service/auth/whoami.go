package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	id, err := uuid.Parse(sessionUUID)
	if err != nil {
		return model.Session{}, model.User{}, errs.ErrSessionNotFound
	}

	session, err := s.sessionRepo.Get(ctx, id)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	user, err := s.userRepo.GetByUUID(ctx, session.UserUUID)
	if err != nil {
		return model.Session{}, model.User{}, fmt.Errorf("получить владельца сессии: %w", err)
	}

	return session, user, nil
}
