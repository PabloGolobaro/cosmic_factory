package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
)

func (s *service) Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), s.bcryptCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("хэширование пароля: %w", err)
	}

	user := model.User{
		UUID:         uuid.New(),
		Login:        in.Login,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return uuid.Nil, err
	}

	return user.UUID, nil
}
