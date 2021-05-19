package jobrun

import (
	"errors"
	"fmt"
)

func serialError(n int, name string, err error) error {
	if name != "" {
		return fmt.Errorf("job:%s (serial job #%d) failed: %w", name, n, err)
	}
	return fmt.Errorf("serial job #%d failed: %w", n, err)
}

type parallelError struct {
	n    int
	name string
	err  error
}

func (pe parallelError) Error() string {
	if pe.name != "" {
		return fmt.Sprintf("job:%s (parallel job #%d) failed: %s", pe.name, pe.n, pe.err)
	}
	return fmt.Sprintf("parallel job #%d failed: %s", pe.n, pe.err)
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

type errorArray []error

func (ea errorArray) Error() string {
	return ea[0].Error()
}

// Unwrap is provided for errors.Unwrap
func (ea errorArray) Unwrap() error {
	return errors.Unwrap(ea[0])
}

// Is is provided for errors.Is
func (ea errorArray) Is(err error) bool {
	for _, e := range ea {
		if ok := errors.Is(e, err); ok {
			return true
		}
	}
	return false
}

// As is provided for errors.As
func (ea errorArray) As(v interface{}) bool {
	for _, e := range ea {
		if ok := errors.As(e, v); ok {
			return true
		}
	}
	return false
}

// Errors returns all wrapped errors in array.
func (ea errorArray) Errors() []error {
	return []error(ea)
}
