package handler

import (
	job2 "WarpQueue/internal/job"
	"WarpQueue/internal/queue"
	"log"
	"time"
)

type Pool struct {
	queue    queue.Queue
	registry *Registry
}

func NewPool(q queue.Queue, r *Registry) *Pool {
	return &Pool{queue: q, registry: r}
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
		p.queue.UpdateStatus(job.ID, job2.StatusRunning)
		if err != nil {
			time.Sleep(200 * time.Millisecond)

		}

		log.Printf("worker-%d processing job %s", id, job.ID)

		handler, err := p.registry.Get(job.ID)
		if err != nil {
			log.Printf("worker-%d handler lookup failed for %s: %v", id, job.ID, err)
			continue
		}

		if err := handler(job); err != nil {
			log.Printf("worker-%d job %s failed: %v", id, job.ID, err)
			p.queue.UpdateStatus(job.ID, job2.StatusFailed)
		} else {
			p.queue.UpdateStatus(job.ID, job2.StatusCompleted)
			log.Printf("worker-%d job %s completed", id, job.ID)
		}
	}
}
