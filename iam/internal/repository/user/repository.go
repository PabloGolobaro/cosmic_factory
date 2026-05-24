package user

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

type Repository interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetByUUID(ctx context.Context, id uuid.UUID) (model.User, error)
}

type stub struct{}

func NewStub() Repository { return &stub{} }

func (*stub) Create(_ context.Context, _ model.User) error {
	return errors.New("not implemented")
}

func (*stub) GetByLogin(_ context.Context, _ string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (*stub) GetByUUID(_ context.Context, _ uuid.UUID) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}
