package job

import (
	"fmt"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func NewStore() *Store {
	return &Store{
		jobs: make(map[string]*Job),
	}
}

func (s *Store) Save(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobCopy := job
	s.jobs[job.ID] = &jobCopy
	return nil
}

func (s *Store) Get(id string) (Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storedJob, ok := s.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("job %s not found", id)
	}

	return *storedJob, nil
}

func (s *Store) Update(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.jobs[job.ID]; !ok {
		return fmt.Errorf("job %s not found", job.ID)
	}

	jobCopy := job
	s.jobs[job.ID] = &jobCopy
	return nil
}

func (s *Store) ListByStatus(status JobStatus) []Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]Job, 0)
	for _, storedJob := range s.jobs {
		if storedJob.Status == status {
			jobs = append(jobs, *storedJob)
		}
	}

	return jobs
}

func (s *Store) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := Stats{
		Total: len(s.jobs),
	}

	for _, storedJob := range s.jobs {
		switch storedJob.Status {
		case StatusPending:
			stats.Pending++
		case StatusRunning:
			stats.Running++
		case StatusRetrying:
			stats.Retrying++
		case StatusCompleted:
			stats.Completed++
		case StatusFailed:
			stats.Failed++
		}
	}

	return stats
}
