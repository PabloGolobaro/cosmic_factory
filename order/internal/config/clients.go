package config

import (
	"net"
	"time"
)

type inventoryClientConfig struct {
	Host         string        `yaml:"host"          env:"INVENTORY_HOST"          env-default:"localhost"`
	Port         string        `yaml:"port"          env:"INVENTORY_PORT"          env-default:"50051"`
	PingInterval time.Duration `yaml:"ping_interval" env:"INVENTORY_PING_INTERVAL" env-default:"10s"`
	PingTimeout  time.Duration `yaml:"ping_timeout"  env:"INVENTORY_PING_TIMEOUT"  env-default:"10s"`
}

func (c *inventoryClientConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

type paymentClientConfig struct {
	Host         string        `yaml:"host"          env:"PAYMENT_HOST"          env-default:"localhost"`
	Port         string        `yaml:"port"          env:"PAYMENT_PORT"          env-default:"50052"`
	PingInterval time.Duration `yaml:"ping_interval" env:"PAYMENT_PING_INTERVAL" env-default:"10s"`
	PingTimeout  time.Duration `yaml:"ping_timeout"  env:"PAYMENT_PING_TIMEOUT"  env-default:"10s"`
}

func (c *paymentClientConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
