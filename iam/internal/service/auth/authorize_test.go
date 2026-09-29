package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
)

type fakeUserRepo struct{ user model.User }

func (f fakeUserRepo) GetByLogin(context.Context, string) (model.User, error) { return f.user, nil }

func (f fakeUserRepo) GetByUUID(context.Context, uuid.UUID) (model.User, error) { return f.user, nil }

type fakeSessionRepo struct{ session *model.Session }

func (f fakeSessionRepo) Save(context.Context, model.Session) error { return nil }

func (f fakeSessionRepo) Delete(context.Context, uuid.UUID) error { return nil }

func (f fakeSessionRepo) Get(context.Context, uuid.UUID) (model.Session, error) {
	if f.session == nil {
		return model.Session{}, errs.ErrSessionNotFound
	}
	return *f.session, nil
}

type fakePolicy struct {
	allowed bool
	err     error
	got     model.AuthzInput
}

func (f *fakePolicy) Allow(_ context.Context, in model.AuthzInput) (bool, error) {
	f.got = in
	return f.allowed, f.err
}

func TestAuthorize(t *testing.T) {
	user := model.User{UUID: uuid.New(), Login: "u", Role: model.RoleClient}
	session := &model.Session{UUID: uuid.New(), UserUUID: user.UUID}
	ownerUUID := uuid.NewString()

	t.Run("passes subject and resource to policy", func(t *testing.T) {
		policy := &fakePolicy{allowed: true}
		svc := New(fakeUserRepo{user}, fakeSessionRepo{session}, policy, time.Hour)

		decision, err := svc.Authorize(context.Background(), input.AuthorizeInput{
			SessionUUID: session.UUID.String(),
			Action:      "order:read",
			OwnerUUID:   ownerUUID,
		})

		require.NoError(t, err)
		require.True(t, decision.Allowed)
		require.Equal(t, user.UUID, decision.User.UUID)
		require.Equal(t, model.AuthzInput{
			Subject:  model.AuthzSubject{UUID: user.UUID.String(), Role: model.RoleClient},
			Action:   "order:read",
			Resource: model.AuthzResource{OwnerUUID: ownerUUID},
		}, policy.got)
	})

	t.Run("denied decision is not an error", func(t *testing.T) {
		svc := New(fakeUserRepo{user}, fakeSessionRepo{session}, &fakePolicy{}, time.Hour)

		decision, err := svc.Authorize(context.Background(), input.AuthorizeInput{
			SessionUUID: session.UUID.String(),
			Action:      "order:pay",
		})

		require.NoError(t, err)
		require.False(t, decision.Allowed)
	})

	t.Run("invalid session", func(t *testing.T) {
		svc := New(fakeUserRepo{user}, fakeSessionRepo{}, &fakePolicy{allowed: true}, time.Hour)

		_, err := svc.Authorize(context.Background(), input.AuthorizeInput{
			SessionUUID: uuid.NewString(),
			Action:      "order:read",
		})

		require.ErrorIs(t, err, errs.ErrSessionNotFound)
	})

	t.Run("policy error is wrapped", func(t *testing.T) {
		policyErr := errors.New("boom")
		svc := New(fakeUserRepo{user}, fakeSessionRepo{session}, &fakePolicy{err: policyErr}, time.Hour)

		_, err := svc.Authorize(context.Background(), input.AuthorizeInput{
			SessionUUID: session.UUID.String(),
			Action:      "order:read",
		})

		require.ErrorIs(t, err, policyErr)
	})
}
