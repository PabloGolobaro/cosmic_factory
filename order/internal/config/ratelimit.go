package config

type rateLimitConfig struct {
	RedisAddress string `yaml:"redis_address"`
	Rate         int    `yaml:"rate"`
	Burst        int    `yaml:"burst"`
}
