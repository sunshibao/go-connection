package connection

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

type RedisConfig struct {
	host     string
	port     string
	password string
	database int
}

var (
	client *redis.Client
)

func Redis() *redis.Client {
	return client
}

func NewRedisConfig(options ...RedisOption) *RedisConfig {
	config := &RedisConfig{}
	for _, option := range options {
		option(config)
	}
	return config
}

type RedisOption func(config *RedisConfig)

func RedisHost(host string) RedisOption {
	return func(config *RedisConfig) {
		config.host = host
	}
}

func RedisPort(port string) RedisOption {
	return func(config *RedisConfig) {
		config.port = port
	}
}

func RedisPassword(password string) RedisOption {
	return func(config *RedisConfig) {
		config.password = password
	}
}

func RedisDatabase(database int) RedisOption {
	return func(config *RedisConfig) {
		config.database = database
	}
}

func (config *RedisConfig) Connect() error {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     config.getAddress(),
			Password: config.password, // no password set
			DB:       config.database, // use default DB
		})
	}
	return nil
}

func (config *RedisConfig) getAddress() string {
	return fmt.Sprintf("%s:%s", config.host, config.port)
}

func (config *RedisConfig) Close() error {
	return client.Close()
}
