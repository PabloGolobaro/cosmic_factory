package config

type loggerConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
}
