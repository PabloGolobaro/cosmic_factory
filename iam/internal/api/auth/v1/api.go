package authv1

import (
	"context"

	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (model.Session, error)
	Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error)
	Logout(ctx context.Context, sessionUUID string) error
}

type API struct {
	authproto.UnimplementedAuthServiceServer
	authSvc AuthService
}

func NewAPI(authSvc AuthService) *API {
	return &API{authSvc: authSvc}
}
