package jobrun

import (
	"context"
	"sync"
)

// Runner provides job logic to run.
type Runner interface {
	Run(ctx context.Context) error
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
func (p Parallel) Run(ctx0 context.Context) error {
	ctx, cancel := context.WithCancel(ctx0)
	defer cancel()

	errs := &ErrorParallel{}

	var wg sync.WaitGroup
	wg.Add(len(p))
	for i0, r0 := range p {
		go func(i int, r Runner) {
			defer wg.Done()
			err := r.Run(ctx)
			if err != nil {
				errs.add(i, err)
				cancel()
			}
		}(i0, r0)
	}
	wg.Wait()

	if len(errs.errs) > 0 {
		return errs
	}
	return nil
}
