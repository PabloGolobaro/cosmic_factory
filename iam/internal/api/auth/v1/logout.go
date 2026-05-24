package authv1

import (
	"context"
	"errors"

	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
)

func (a *API) Logout(ctx context.Context, req *authproto.LogoutRequest) (*authproto.LogoutResponse, error) {
	if req.GetSessionUuid() == "" {
		return nil, errs.ErrEmptySessionID
	}

	err := a.authSvc.Logout(ctx, req.GetSessionUuid())
	if errors.Is(err, errs.ErrSessionNotFound) {
		// Идемпотентность: повторный Logout не является ошибкой
		return &authproto.LogoutResponse{}, nil
	}

	if err != nil {
		return nil, err
	}

	return &authproto.LogoutResponse{}, nil
}
