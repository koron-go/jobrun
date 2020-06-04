package jobrun

import (
	"context"
	"errors"
	"sync"
)

// Runner provides job logic to run.
type Runner interface {
	Run(ctx context.Context) error
}

// RunFunc is Runner wrapper for function.
type RunFunc func(ctx context.Context) error

// Run implements Runner interface.
func (f RunFunc) Run(ctx context.Context) error {
	return f(ctx)
}

// Serial defines serial job runner.
type Serial []Runner

// Add adds a runner to Serial.
func (s *Serial) Add(r ...Runner) *Serial {
	*s = append(*s, r...)
	return s
}

// Run implements Runner.
func (s Serial) Run(ctx context.Context) error {
	for i, r := range s {
		err := r.Run(ctx)
		if err != nil {
			return serialError(i, err)
		}
	}
	return nil
}

// Parallel defines parallel job runner.
type Parallel []Runner

// Add adds a runner to Parallel.
func (p *Parallel) Add(r ...Runner) *Parallel {
	*p = append(*p, r...)
	return p
}

// Run implements Runner.
func (p Parallel) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var mu sync.Mutex
	errs := make(errorArray, 0, len(p))

	var wg sync.WaitGroup
	wg.Add(len(p))
	for i, job := range p {
		go func(n int, r Runner) {
			defer wg.Done()
			err := r.Run(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				mu.Lock()
				errs = append(errs, parallelError{n: n, err: err})
				mu.Unlock()
				cancel()
			}
		}(i, job)
	}
	wg.Wait()

	if len(errs) > 0 {
		return errs
	}
	return nil
}
