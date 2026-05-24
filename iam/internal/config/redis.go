package config

import "net"

type redisConfig struct {
	Host string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
}

func (c *redisConfig) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}
