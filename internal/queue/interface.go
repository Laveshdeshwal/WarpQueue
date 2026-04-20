package queue

import (
	"WarpQueue/internal/job"
	"errors"
)

type Queue interface {
	Enqueue(job job.Job) error
	Dequeue() (job.Job, error)
	Size() int
	UpdateStatus(id string, status job.JobStatus) error
	GetJob(id string) (job.Job, error)
}

var ErrQueueEmpty = errors.New("queue is empty")
