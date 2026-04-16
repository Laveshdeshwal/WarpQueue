package worker

import (
	"WarpQueue/internal/queue"
	handler "WarpQueue/internal/worker"
	"log"
)

func main() {
	q := queue.NewMemoryQueue()
	r := handler.NewRegistry()

	// Example handler
	r.Register("send_email", func(job handler.Job) error {
		log.Printf("sending email for job %s", job.ID)
		return nil
	})

	pool := handler.NewPool(q, r)
	pool.Start(3)

	log.Println("worker pool started with 3 workers")
	select {}
}
