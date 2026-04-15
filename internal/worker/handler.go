package handler

import (
	"WarpQueue/internal/job"
	"fmt"
	"sync"
)

type Job = job.Job
type Handler func(Job) error

type Registry struct {
	handlers map[string]Handler
	mu       sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]Handler)}
}

func (r *Registry) Register(name string, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[name] = handler
}
func (r *Registry) Get(name string) (Handler, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[name]
	if !ok {
		return nil, fmt.Errorf("Handler %s not found", name)
	}
	return h, nil

}
