package config

import (
	"net"
	"time"
)

type httpConfig struct {
	Host              string        `yaml:"host"               env:"HTTP_HOST"               env-default:"localhost"`
	Port              string        `yaml:"port"               env:"HTTP_PORT"               env-default:"8080"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env:"HTTP_READ_HEADER_TIMEOUT" env-default:"5s"`
	ReadTimeout       time.Duration `yaml:"read_timeout"       env:"HTTP_READ_TIMEOUT"       env-default:"15s"`
	WriteTimeout      time.Duration `yaml:"write_timeout"      env:"HTTP_WRITE_TIMEOUT"      env-default:"15s"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"       env:"HTTP_IDLE_TIMEOUT"       env-default:"60s"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"   env:"HTTP_SHUTDOWN_TIMEOUT"   env-default:"10s"`
}

func (c *httpConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
