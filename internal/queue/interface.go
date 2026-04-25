package queue

import (
	"WarpQueue/internal/job"
	"errors"
)

type Queue interface {
	Enqueue(job job.Job) error
	Dequeue() (job.Job, error)
	Size() int
	Save(job job.Job) error
	Get(id string) (job.Job, error)
	Update(job job.Job) error
	ListByStatus(status job.JobStatus) []job.Job
	Stats() job.Stats
}

var ErrQueueEmpty = errors.New("queue is empty")
