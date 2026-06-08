package ratelimit

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

// Middleware создаёт HTTP middleware с распределённым rate limiter.
// Ключ — путь запроса (r.URL.Path), лимит общий для всех инстансов через Redis.
// Fail-open: если Redis недоступен — запрос пропускается дальше.
func Middleware(limiter *redis_rate.Limiter, rate, burst int) func(http.Handler) http.Handler {
	limit := redis_rate.Limit{
		Rate:   rate,
		Burst:  burst,
		Period: time.Second,
	}

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
