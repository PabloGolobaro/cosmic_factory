package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
)

type Service interface {
	Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}

type service struct {
	userRepo   userRepository
	bcryptCost int
}

func New(userRepo userRepository, bcryptCost int) Service {
	return &service{userRepo: userRepo, bcryptCost: bcryptCost}
}
