package authv1

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/converter"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
)

func (a *API) Authorize(ctx context.Context, req *authproto.AuthorizeRequest) (*authproto.AuthorizeResponse, error) {
	decision, err := a.authSvc.Authorize(ctx, input.AuthorizeInput{
		SessionUUID: req.GetSessionUuid(),
		Action:      req.GetAction(),
		OwnerUUID:   req.GetResource().GetOwnerUuid(),
	})
	if err != nil {
		return nil, err
	}

	return &authproto.AuthorizeResponse{
		Allowed:  decision.Allowed,
		UserUuid: decision.User.UUID.String(),
		Role:     converter.RoleToProto(decision.User.Role),
	}, nil
}
