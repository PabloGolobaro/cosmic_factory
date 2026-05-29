package config

type otelConfig struct {
	Endpoint    string `yaml:"endpoint"     env:"OTEL_ENDPOINT"     env-default:"localhost:4317"`
	ServiceName string `yaml:"service_name" env:"OTEL_SERVICE_NAME" env-default:"payment-service"`
}
