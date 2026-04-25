package queue

import (
	"WarpQueue/internal/job"
	"sync"
)

type MemoryQueue struct {
	mu    sync.RWMutex
	items []string
	store *job.Store
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		items: make([]string, 0),
		store: job.NewStore(),
	}
}

func (q *MemoryQueue) Save(newJob job.Job) error {
	return q.store.Save(newJob)
}

func (q *MemoryQueue) Get(id string) (job.Job, error) {
	return q.store.Get(id)
}

func (q *MemoryQueue) Update(updatedJob job.Job) error {
	return q.store.Update(updatedJob)
}

func (q *MemoryQueue) ListByStatus(status job.JobStatus) []job.Job {
	return q.store.ListByStatus(status)
}

func (q *MemoryQueue) Stats() job.Stats {
	return q.store.Stats()
}

func (q *MemoryQueue) Enqueue(newJob job.Job) error {
	newJob.Status = job.StatusPending
	if err := q.store.Save(newJob); err != nil {
		return err
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, newJob.ID)
	return nil
}

func (q *MemoryQueue) Dequeue() (job.Job, error) {
	q.mu.Lock()
	if len(q.items) == 0 {
		q.mu.Unlock()
		return job.Job{}, ErrQueueEmpty
	}
	jobID := q.items[0]
	q.items = q.items[1:]
	q.mu.Unlock()

	storedJob, err := q.store.Get(jobID)
	if err != nil {
		return job.Job{}, err
	}

	storedJob.Attempts++
	storedJob.Status = job.StatusRunning
	if err := q.store.Update(storedJob); err != nil {
		return job.Job{}, err
	}

	return storedJob, nil
}

func (q *MemoryQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items)
}
