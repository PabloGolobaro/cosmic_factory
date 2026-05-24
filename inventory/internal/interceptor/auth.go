package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/auth"
)

type iamClient interface {
	Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error)
}

func New(client iamClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "отсутствует metadata")
		}

		values := md.Get("session-uuid")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "отсутствует session-uuid")
		}

		userUUID, err := client.Whoami(ctx, values[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "недействительная сессия")
		}

		return handler(auth.WithUserUUID(ctx, userUUID), req)
	}
}
