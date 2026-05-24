package user

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

func (s *store) GetByLogin(ctx context.Context, login string) (model.User, error) {
	return s.queryUser(ctx, selectUserByLogin, login)
}
