package queue

import (
	"WarpQueue/internal/job"
	"sync"
)

type MemoryQueue struct {
	mu   sync.RWMutex
	item []job.Job
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{item: make([]job.Job, 0)}
}
func (q *MemoryQueue) Enqueue(job job.Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.item = append(q.item, job)
	return nil
}
func (q *MemoryQueue) Dequeue() (job.Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.item) == 0 {
		return job.Job{}, ErrQueueEmpty
	}
	job := q.item[0]
	q.item = q.item[1:]
	return job, nil
}
func (q *MemoryQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.item)
}
