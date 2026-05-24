package userv1

import (
	"context"

	userproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/user/v1"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/service/input"
)

type UserService interface {
	Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}

type API struct {
	userproto.UnimplementedUserServiceServer
	userSvc UserService
}

func NewAPI(userSvc UserService) *API {
	return &API{userSvc: userSvc}
}
