package queue

import (
	"WarpQueue/internal/job"
	"errors"
)

type Queue interface {
	Enqueue(job job.Job) error
	Dequeue() (job.Job, error)
	Size() int
}

var ErrQueueEmpty = errors.New("queue is empty")
