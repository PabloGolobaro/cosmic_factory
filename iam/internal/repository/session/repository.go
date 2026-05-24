package session

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

type Repository interface {
	Save(ctx context.Context, session model.Session) error
	Get(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error)
	Delete(ctx context.Context, sessionUUID uuid.UUID) error
}

type stub struct{}

func NewStub() Repository { return &stub{} }

func (*stub) Save(_ context.Context, _ model.Session) error {
	return errors.New("not implemented")
}

func (*stub) Get(_ context.Context, _ uuid.UUID) (model.Session, error) {
	return model.Session{}, errors.New("not implemented")
}

func (*stub) Delete(_ context.Context, _ uuid.UUID) error {
	return errors.New("not implemented")
}
