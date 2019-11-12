package jobrun

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

func serialError(n int, err error) error {
	return fmt.Errorf("serial job #%d failed: %w", n, err)
}

type parallelError struct {
	n   int
	err error
}

func (pe parallelError) Error() string {
	return fmt.Sprintf("parallel job #%d failed: %s", pe.err)
}

func (pe parallelError) Unwrap() error {
	return pe.err
}

func (pe parallelError) Is(err error) bool {
	return errors.Is(pe.err, err)
}

func (pe parallelError) As(v interface{}) bool {
	return errors.As(pe.err, v)
}

// ErrorParallel is set of errors which occur in parallel.
type ErrorParallel struct {
	mu   sync.Mutex
	errs []parallelError
}

func (e *ErrorParallel) add(n int, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	e.mu.Lock()
	e.errs = append(e.errs, parallelError{n: n, err: err})
	e.mu.Unlock()
}

func (e *ErrorParallel) Error() string {
	return e.errs[0].Error()
}

// Unwrap is provided for errors.Unwrap
func (e *ErrorParallel) Unwrap() error {
	return e.errs[0]
}

// Is is provided for errors.Is
func (e *ErrorParallel) Is(err error) bool {
	for _, pe := range e.errs {
		if ok := pe.Is(err); ok {
			return true
		}
	}
	return false
}

// As is provided for errors.As
func (e *ErrorParallel) As(v interface{}) bool {
	for _, pe := range e.errs {
		if ok := pe.As(v); ok {
			return true
		}
	}
	return false
}

// Errors returns all wrapped errors in array.
func (e *ErrorParallel) Errors() []error {
	errs := make([]error, len(e.errs))
	for i, pe := range e.errs {
		errs[i] = pe
	}
	return errs
}
