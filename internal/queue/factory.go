package queue

import (
	"WarpQueue/internal/config"
	"fmt"
)

func NewFromConfig(cfg config.Config) (Queue, error) {
	switch cfg.QueueType {
	case "memory":
		return NewMemoryQueue(), nil
	case "redis":
		return NewRedisQueue(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.RedisKeyPrefix)
	default:
		return nil, fmt.Errorf("unsupported queue type: %s", cfg.QueueType)
	}
}
