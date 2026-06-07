package ratelimit

import (
	"log/slog"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

// NewHTTPMiddleware создаёт HTTP middleware с распределённым rate limiter.
// Ключ — путь запроса (r.URL.Path), лимит общий для всех инстансов через Redis.
// Fail-open: если Redis недоступен — запрос пропускается дальше.
func NewHTTPMiddleware(limiter *redis_rate.Limiter, limit redis_rate.Limit) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Path

			res, err := limiter.Allow(r.Context(), key, limit)
			if err != nil {
				slog.Error("ошибка проверки rate limit", "path", key, "error", err)
				next.ServeHTTP(w, r)
				return
			}

			if res.Allowed == 0 {
				slog.Warn("запрос отклонён rate limiter", "path", key, "retry_after", res.RetryAfter)
				http.Error(w, "слишком много запросов, попробуйте позже", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
