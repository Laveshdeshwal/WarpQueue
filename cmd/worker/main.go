package main

import (
	"WarpQueue/internal/config"
	"WarpQueue/internal/logger"
	"WarpQueue/internal/queue"
	handler "WarpQueue/internal/worker"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()
	logger.InitLogger("worker-log", cfg.LogLevel)

	q, err := queue.NewFromConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	r := handler.NewRegistry()

	// Example handler
	r.Register("send_email", func(job handler.Job) error {
		log.Printf("sending email for job %s", job.ID)
		return nil
	})

	pool := handler.NewPool(q, r)
	pool.Start(cfg.WorkerCount)

	log.Printf("worker pool started with %d workers", cfg.WorkerCount)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Printf("worker shutting down with timeout %s", cfg.ShutdownTimeout)
}
