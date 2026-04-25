package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultServerPort      = "8080"
	defaultWorkerCount     = 3
	defaultLogLevel        = "info"
	defaultQueueType       = "memory"
	defaultShutdownTimeout = 10 * time.Second
)

type Config struct {
	ServerPort      string
	WorkerCount     int
	LogLevel        string
	QueueType       string
	ShutdownTimeout time.Duration
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	RedisKeyPrefix  string
}

func Load() Config {
	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("SERVER_PORT", defaultServerPort)
	v.SetDefault("WORKER_COUNT", defaultWorkerCount)
	v.SetDefault("LOG_LEVEL", defaultLogLevel)
	v.SetDefault("QUEUE_TYPE", defaultQueueType)
	v.SetDefault("SHUTDOWN_TIMEOUT", defaultShutdownTimeout.String())
	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("REDIS_KEY_PREFIX", "warpqueue")

	return Config{
		ServerPort:      v.GetString("SERVER_PORT"),
		WorkerCount:     v.GetInt("WORKER_COUNT"),
		LogLevel:        strings.ToLower(v.GetString("LOG_LEVEL")),
		QueueType:       strings.ToLower(v.GetString("QUEUE_TYPE")),
		ShutdownTimeout: v.GetDuration("SHUTDOWN_TIMEOUT"),
		RedisAddr:       v.GetString("REDIS_ADDR"),
		RedisPassword:   v.GetString("REDIS_PASSWORD"),
		RedisDB:         v.GetInt("REDIS_DB"),
		RedisKeyPrefix:  v.GetString("REDIS_KEY_PREFIX"),
	}
}
