package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/PabloGolobaro/cosmic_factory/iam/internal/errors"
)

// ErrorInterceptor конвертирует доменные ошибки IAM в gRPC-статусы.
// Уже упакованные gRPC-статусы (protovalidate, RecoveryInterceptor) пропускаются без изменений.
// Все неизвестные ошибки маппируются в codes.Internal без раскрытия деталей.
func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		return nil, domainToGRPC(err)
	}
}

func domainToGRPC(err error) error {
	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, errs.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, errs.ErrInvalidLogin),
		errors.Is(err, errs.ErrWeakPassword),
		errors.Is(err, errs.ErrEmptyCredential),
		errors.Is(err, errs.ErrEmptySessionID),
		errors.Is(err, errs.ErrInvalidUUID):
		return status.Error(codes.InvalidArgument, err.Error())

	default:
		return status.Error(codes.Internal, "внутренняя ошибка сервера")
	}
}
