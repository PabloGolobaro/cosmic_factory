package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/auth"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

const SessionMetadataKey = "session-uuid"

// publicMethods — методы, не требующие аутентификации (health-check).
var publicMethods = map[string]struct{}{
	grpc_health_v1.Health_Check_FullMethodName:              {},
	grpc_health_v1.Health_List_FullMethodName:               {},
	inventoryv1.InventoryService_CommitParts_FullMethodName: {},
}

type iamClient interface {
	Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error)
}

// Auth создаёт gRPC-interceptor аутентификации через IAM.
func Auth(client iamClient) grpc.UnaryServerInterceptor { return New(client) }

func New(client iamClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := publicMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "отсутствует metadata")
		}

		values := md.Get(SessionMetadataKey)
		if len(values) == 0 || values[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "отсутствует session-uuid")
		}

		userUUID, err := client.Whoami(ctx, values[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "недействительная сессия")
		}

		return handler(auth.WithUserUUID(ctx, userUUID), req)
	}
}
