package authv1

import (
	"context"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
)

func (a *API) Login(ctx context.Context, req *authproto.LoginRequest) (*authproto.LoginResponse, error) {
	if req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, errs.ErrEmptyCredential
	}

	session, err := a.authSvc.Login(ctx, input.LoginInput{Login: req.GetLogin(), Password: req.GetPassword()})
	if err != nil {
		return nil, err
	}

	return &authproto.LoginResponse{SessionUuid: session.UUID.String()}, nil
}
