package authv1

import (
	"context"

	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
)

func (a *API) Login(ctx context.Context, req *authproto.LoginRequest) (*authproto.LoginResponse, error) {
	if req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, errs.ErrEmptyCredential
	}

	session, err := a.authSvc.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &authproto.LoginResponse{SessionUuid: session.UUID.String()}, nil
}
