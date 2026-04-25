package main

import (
	"WarpQueue/internal/api"
	"WarpQueue/internal/config"
	"WarpQueue/internal/logger"
	"WarpQueue/internal/queue"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()
	logger.InitLogger("app-log", cfg.LogLevel)

	q, err := newQueue(cfg.QueueType)
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewServer(q, cfg.WorkerCount)
	httpServer := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: server.Routes(),
	}

	logger.Info("Server running on :" + cfg.ServerPort)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func newQueue(queueType string) (queue.Queue, error) {
	switch queueType {
	case "memory":
		return queue.NewMemoryQueue(), nil
	default:
		return nil, fmt.Errorf("unsupported queue type: %s", queueType)
	}
}
