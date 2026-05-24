package userv1

import (
	"context"

	userproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/user/v1"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
)

func (a *API) Register(ctx context.Context, req *userproto.RegisterRequest) (*userproto.RegisterResponse, error) {
	login := req.GetInfo().GetInfo().GetLogin()
	password := req.GetInfo().GetPassword()

	if login == "" {
		return nil, errs.ErrInvalidLogin
	}

	if len(password) < 8 {
		return nil, errs.ErrWeakPassword
	}

	userUUID, err := a.userSvc.Register(ctx, input.RegisterInput{Login: login, Password: password})
	if err != nil {
		return nil, err
	}

	return &userproto.RegisterResponse{UserUuid: userUUID.String()}, nil
}
