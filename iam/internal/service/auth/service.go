package auth

import (
	"context"
	"errors"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

type Service interface {
	Login(ctx context.Context, login, password string) (model.Session, error)
	Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error)
	Logout(ctx context.Context, sessionUUID string) error
}

type stub struct{}

func NewStub() Service { return &stub{} }

func (*stub) Login(_ context.Context, _, _ string) (model.Session, error) {
	return model.Session{}, errors.New("not implemented")
}

func (*stub) Whoami(_ context.Context, _ string) (model.Session, model.User, error) {
	return model.Session{}, model.User{}, errors.New("not implemented")
}

func (*stub) Logout(_ context.Context, _ string) error {
	return errors.New("not implemented")
}
