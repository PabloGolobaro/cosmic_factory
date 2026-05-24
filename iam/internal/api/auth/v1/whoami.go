package authv1

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/converter"
	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
)

func (a *API) Whoami(ctx context.Context, req *authproto.WhoamiRequest) (*authproto.WhoamiResponse, error) {
	if req.GetSessionUuid() == "" {
		return nil, errs.ErrEmptySessionID
	}

	session, user, err := a.authSvc.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		return nil, err
	}

	return &authproto.WhoamiResponse{
		Session: converter.SessionToProto(session),
		User:    converter.UserToProto(user),
	}, nil
}
