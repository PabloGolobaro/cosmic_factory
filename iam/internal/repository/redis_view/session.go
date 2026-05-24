package redis_view

// SessionRedisView — представление сессии для хранения в Redis как HashMap.
// Поля соответствуют ключам хэшмапы (теги redis:"...").
// UpdatedAt не хранится: сессии не обновляются, только создаются и удаляются.
type SessionRedisView struct {
	UUID      string `redis:"uuid"`
	UserUUID  string `redis:"user_uuid"`
	Login     string `redis:"login"`
	CreatedAt string `redis:"created_at"`
	ExpiresAt string `redis:"expires_at"`
}
