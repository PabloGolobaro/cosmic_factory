package v1

import (
	"context"

	"github.com/google/uuid"

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
