package config

type kafkaConfig struct {
	Brokers       []string `yaml:"brokers"        env:"KAFKA_BROKERS"        env-separator:","`
	ConsumerGroup string   `yaml:"consumer_group" env:"KAFKA_CONSUMER_GROUP" env-default:"order-service"`
	ConsumeTopic  string   `yaml:"consume_topic"  env:"KAFKA_CONSUME_TOPIC"  env-default:"assembly.ship-assembled"`
	ProduceTopic  string   `yaml:"produce_topic"  env:"KAFKA_PRODUCE_TOPIC"  env-default:"order.paid"`
}
