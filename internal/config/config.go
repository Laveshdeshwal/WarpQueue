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
}

func Load() Config {
	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("SERVER_PORT", defaultServerPort)
	v.SetDefault("WORKER_COUNT", defaultWorkerCount)
	v.SetDefault("LOG_LEVEL", defaultLogLevel)
	v.SetDefault("QUEUE_TYPE", defaultQueueType)
	v.SetDefault("SHUTDOWN_TIMEOUT", defaultShutdownTimeout.String())

	return Config{
		ServerPort:      v.GetString("SERVER_PORT"),
		WorkerCount:     v.GetInt("WORKER_COUNT"),
		LogLevel:        strings.ToLower(v.GetString("LOG_LEVEL")),
		QueueType:       strings.ToLower(v.GetString("QUEUE_TYPE")),
		ShutdownTimeout: v.GetDuration("SHUTDOWN_TIMEOUT"),
	}
}
