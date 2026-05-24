package config

import (
	"net"
	"time"
)

type iamClientConfig struct {
	Host         string        `yaml:"host"          env:"IAM_HOST"          env-default:"localhost"`
	Port         string        `yaml:"port"          env:"IAM_PORT"          env-default:"50053"`
	PingInterval time.Duration `yaml:"ping_interval" env:"IAM_PING_INTERVAL" env-default:"10s"`
	PingTimeout  time.Duration `yaml:"ping_timeout"  env:"IAM_PING_TIMEOUT"  env-default:"10s"`
}

func (c *iamClientConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
