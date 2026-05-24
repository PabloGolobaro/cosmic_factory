package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
)

func (s *service) Login(ctx context.Context, in input.LoginInput) (model.Session, error) {
	user, err := s.userRepo.GetByLogin(ctx, in.Login)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return model.Session{}, errs.ErrInvalidCredentials
		}

		return model.Session{}, fmt.Errorf("получить пользователя: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return model.Session{}, errs.ErrInvalidCredentials
	}

	now := time.Now()
	session := model.Session{
		UUID:      uuid.New(),
		UserUUID:  user.UUID,
		Login:     user.Login,
		CreatedAt: now,
		ExpiresAt: now.Add(s.ttl),
	}

	if err = s.sessionRepo.Save(ctx, session); err != nil {
		return model.Session{}, fmt.Errorf("сохранить сессию: %w", err)
	}

	return session, nil
}
