package userv1

import (
	"context"

	"github.com/google/uuid"
	userproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/user/v1"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/converter"
	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
)

func (a *API) GetUser(ctx context.Context, req *userproto.GetUserRequest) (*userproto.GetUserResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, errs.ErrInvalidUUID
	}

	user, err := a.userSvc.GetUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return &userproto.GetUserResponse{User: converter.UserToProto(user)}, nil
}
