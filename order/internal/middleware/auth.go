package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/auth"
)

type iamClient interface {
	Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error)
}

func New(client iamClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "отсутствует заголовок Authorization", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "неверный формат Authorization", http.StatusUnauthorized)
				return
			}

			sessionUUID := parts[1]

			userUUID, err := client.Whoami(r.Context(), sessionUUID)
			if err != nil {
				http.Error(w, "недействительная сессия", http.StatusUnauthorized)
				return
			}

			ctx := auth.WithUserUUID(r.Context(), userUUID)
			ctx = auth.WithSessionUUID(ctx, sessionUUID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
