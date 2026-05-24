package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/auth"
)

func New() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if sessionUUID, ok := auth.SessionUUIDFromContext(ctx); ok {
			ctx = metadata.AppendToOutgoingContext(ctx, "session-uuid", sessionUUID)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
