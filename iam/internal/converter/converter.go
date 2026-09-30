package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	commonv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/common/v1"
)

func UserToProto(u model.User) *commonv1.User {
	pb := &commonv1.User{
		Uuid:      u.UUID.String(),
		Info:      &commonv1.UserInfo{Login: u.Login},
		CreatedAt: timestamppb.New(u.CreatedAt),
		Role:      RoleToProto(u.Role),
	}

	if u.UpdatedAt != nil {
		pb.UpdatedAt = timestamppb.New(*u.UpdatedAt)
	}

	return pb
}

func SessionToProto(s model.Session) *commonv1.Session {
	pb := &commonv1.Session{
		Uuid:      s.UUID.String(),
		CreatedAt: timestamppb.New(s.CreatedAt),
		ExpiresAt: timestamppb.New(s.ExpiresAt),
	}

	if s.UpdatedAt != nil {
		pb.UpdatedAt = timestamppb.New(*s.UpdatedAt)
	}

	return pb
}

func RoleToProto(r model.Role) commonv1.Role {
	switch r {
	case model.RoleClient:
		return commonv1.Role_ROLE_CLIENT
	case model.RoleManager:
		return commonv1.Role_ROLE_MANAGER
	default:
		return commonv1.Role_ROLE_UNSPECIFIED
	}
}
