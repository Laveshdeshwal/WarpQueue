package handler

import (
	job2 "WarpQueue/internal/job"
	"WarpQueue/internal/queue"
	"errors"
	"log"
	"time"
)

type Pool struct {
	queue      queue.Queue
	registry   *Registry
	retryDelay time.Duration
}

func NewPool(q queue.Queue, r *Registry) *Pool {
	return &Pool{
		queue:      q,
		registry:   r,
		retryDelay: time.Second,
	}
}

func (p *Pool) SetRetryDelay(delay time.Duration) {
	p.retryDelay = delay
}

func (p *Pool) Start(n int) {
	for i := 0; i < n; i++ {
		workerId := i
		go p.run(workerId)
	}
}

func (p *Pool) run(id int) {
	for {
		job, err := p.queue.Dequeue()
		if err != nil {
			if errors.Is(err, queue.ErrQueueEmpty) {
				time.Sleep(200 * time.Millisecond)
				continue
			}

			log.Printf("worker-%d dequeue failed: %v", id, err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		log.Printf("worker-%d processing job %s", id, job.ID)

		handler, err := p.registry.Get(job.Type)
		if err != nil {
			log.Printf("worker-%d handler lookup failed for %s: %v", id, job.ID, err)
			p.handleFailure(job, err)
			continue
		}

		if err := handler(job); err != nil {
			log.Printf("worker-%d job %s failed: %v", id, job.ID, err)
			p.handleFailure(job, err)
		} else {
			job.Status = job2.StatusCompleted
			if err := p.queue.Update(job); err != nil {
				log.Printf("worker-%d job %s status update failed: %v", id, job.ID, err)
				continue
			}
			log.Printf("worker-%d job %s completed", id, job.ID)
		}
	}
}

func (p *Pool) handleFailure(job job2.Job, handlerErr error) {
	if job.Attempts < job.MaxRetries {
		job.Status = job2.StatusRetrying
		if err := p.queue.Update(job); err != nil {
			log.Printf("job %s retry status update failed: %v", job.ID, err)
			return
		}

		log.Printf("job %s scheduled for retry %d/%d", job.ID, job.Attempts, job.MaxRetries)
		time.AfterFunc(p.retryDelay, func() {
			if err := p.queue.Enqueue(job); err != nil {
				log.Printf("job %s requeue failed: %v", job.ID, err)
			}
		})
		return
	}

	job.Status = job2.StatusFailed
	if err := p.queue.Update(job); err != nil {
		log.Printf("job %s final failure status update failed: %v", job.ID, err)
		return
	}

	log.Printf("job %s exhausted retries after error: %v", job.ID, handlerErr)
}
