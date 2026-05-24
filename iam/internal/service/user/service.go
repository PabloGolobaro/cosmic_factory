package user

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

type Service interface {
	Register(ctx context.Context, login, password string) (uuid.UUID, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}

type stub struct{}

func NewStub() Service { return &stub{} }

func (*stub) Register(_ context.Context, _, _ string) (uuid.UUID, error) {
	return uuid.UUID{}, errors.New("not implemented")
}

func (*stub) GetUser(_ context.Context, _ uuid.UUID) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}
