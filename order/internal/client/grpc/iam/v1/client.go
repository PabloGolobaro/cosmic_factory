package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/auth"
	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
)

type Client struct {
	inner authproto.AuthServiceClient
}

func New(c authproto.AuthServiceClient) *Client {
	return &Client{inner: c}
}

func (c *Client) Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error) {
	resp, err := c.inner.Whoami(ctx, &authproto.WhoamiRequest{SessionUuid: sessionUUID})
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(resp.GetUser().GetUuid())
}

// Authorize asks IAM whether the session owner may perform action on a resource owned by ownerUUID.
// Pass uuid.Nil for actions without a concrete resource (e.g. create).
func (c *Client) Authorize(ctx context.Context, action string, ownerUUID uuid.UUID) error {
	sessionUUID, ok := auth.SessionUUIDFromContext(ctx)
	if !ok {
		return errs.ErrUnauthorized
	}

	req := &authproto.AuthorizeRequest{SessionUuid: sessionUUID, Action: action}
	if ownerUUID != uuid.Nil {
		req.Resource = &authproto.Resource{OwnerUuid: ownerUUID.String()}
	}

	resp, err := c.inner.Authorize(ctx, req)
	if err != nil {
		return fmt.Errorf("авторизация в IAM: %w", err)
	}

	if !resp.GetAllowed() {
		return fmt.Errorf("%w: %s", errs.ErrForbidden, action)
	}

	return nil
}
