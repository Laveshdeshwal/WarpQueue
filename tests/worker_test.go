package tests

import (
	"WarpQueue/internal/job"
	"WarpQueue/internal/queue"
	handler "WarpQueue/internal/worker"
	"errors"
	"testing"
	"time"
)

func TestPoolRetriesFailedJobs(t *testing.T) {
	q := queue.NewMemoryQueue()
	registry := handler.NewRegistry()
	attempts := 0

	registry.Register("send_email", func(job handler.Job) error {
		attempts++
		if attempts == 1 {
			return errors.New("temporary failure")
		}
		return nil
	})

	pool := handler.NewPool(q, registry)
	pool.SetRetryDelay(10 * time.Millisecond)
	pool.Start(1)

	if err := q.Enqueue(job.Job{
		ID:         "job-1",
		Type:       "send_email",
		MaxRetries: 2,
	}); err != nil {
		t.Fatalf("enqueue job: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		storedJob, err := q.Get("job-1")
		if err != nil {
			t.Fatalf("get job: %v", err)
		}

		if storedJob.Status == job.StatusCompleted {
			if storedJob.Attempts != 2 {
				t.Fatalf("expected 2 attempts, got %d", storedJob.Attempts)
			}
			return
		}

		time.Sleep(20 * time.Millisecond)
	}

	t.Fatal("job did not complete after retry")
}
