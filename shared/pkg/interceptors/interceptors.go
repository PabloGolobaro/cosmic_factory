package interceptors

import (
	"context"
	"log/slog"
	"path"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggerInterceptor создает серверный унарный интерцептор, который логирует
// информацию о времени выполнения методов gRPC сервера.
func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Извлекаем имя метода из полного пути
		method := path.Base(info.FullMethod)

		// Логируем начало вызова метода
		slog.Info("🚀 начало gRPC метода", "method", method)

		// Засекаем время начала выполнения
		startTime := time.Now()

		// Вызываем обработчик
		resp, err := handler(ctx, req)

		// Вычисляем длительность выполнения
		duration := time.Since(startTime)

		// Форматируем сообщение в зависимости от результата
		if err != nil {
			st, _ := status.FromError(err)
			slog.Error("❌ gRPC метод завершён с ошибкой", "method", method, "code", st.Code(), "error", err, "duration", duration)
		} else {
			slog.Info("✅ gRPC метод завершён успешно", "method", method, "duration", duration)
		}

		return resp, err
	}
}

// RecoveryInterceptor создает серверный унарный интерцептор, который перехватывает
// паники в обработчиках и конвертирует их в gRPC ошибку Internal.
// Без этого интерцептора паника в обработчике уронит весь сервер.
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				method := path.Base(info.FullMethod)
				slog.Error("🔥 паника в gRPC методе",
					"method", method,
					"panic", r,
					"stack", string(debug.Stack()),
				)

				err = status.Errorf(codes.Internal, "внутренняя ошибка сервера")
			}
		}()

		return handler(ctx, req)
	}
}
