package queue

import (
	"WarpQueue/internal/job"
	"sync"
)

type MemoryQueue struct {
	mu   sync.RWMutex
	item []job.Job
	jobs map[string]job.Job
}

func (q *MemoryQueue) UpdateStatus(id string, status job.JobStatus) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	job, ok := q.jobs[id]
	if !ok {
		return ErrQueueEmpty
	}
	job.Status = status
	q.jobs[id] = job
	return nil
}

func (q *MemoryQueue) GetJob(id string) (job.Job, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	job, _ := q.jobs[id]
	return job, nil
}
func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{item: make([]job.Job, 0), jobs: make(map[string]job.Job)}
}
func (q *MemoryQueue) Enqueue(job job.Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.item = append(q.item, job)
	q.jobs[job.ID] = job
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
