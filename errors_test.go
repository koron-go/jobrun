package jobrun

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestSerialError(t *testing.T) {
	// normal job
	err := serialError(123, "", fmt.Errorf("hello world"))
	if err.Error() != "serial job #123 failed: hello world" {
		t.Fatal("unexpected Error()")
	}

	// named job
	err = serialError(123, "foobar", fmt.Errorf("hello world"))
	if got := err.Error(); got != "job:foobar (serial job #123) failed: hello world" {
		t.Fatalf("unexpected Error(): got=%s", got)
	}
}

func TestParallelErrors(t *testing.T) {
	// normal job
	err := parallelError{
		n:   123,
		err: fmt.Errorf("hello world"),
	}
	if err.Error() != "parallel job #123 failed: hello world" {
		t.Fatal("unexpected Error()")
	}

	// named job
	err = parallelError{
		n:    123,
		name: "foobar",
		err:  fmt.Errorf("hello world"),
	}
	if got := err.Error(); got != "job:foobar (parallel job #123) failed: hello world" {
		t.Fatalf("unexpected Error(): got=%s", got)
	}
}

func TestParallelErrors_Unwrap(t *testing.T) {
	err := parallelError{err: io.EOF}
	if x := errors.Unwrap(err); x != io.EOF {
		t.Fatalf("errors.Unwrap failed: got=%s", x)
	}
}

func TestParallelErrors_Is(t *testing.T) {
	err := parallelError{err: io.EOF}
	if !errors.Is(err, io.EOF) {
		t.Fatal("errors.Is failed")
	}
}

func TestParallelErrors_As(t *testing.T) {
	err := parallelError{err: &os.PathError{
		Op:   "dummyOp",
		Path: "dummyPath",
		Err:  io.EOF,
	}}
	var xerr *os.PathError
	if !errors.As(err, &xerr) {
		t.Fatal("errors.As failed")
	}
	if !reflect.DeepEqual(xerr, &os.PathError{
		Op:   "dummyOp",
		Path: "dummyPath",
		Err:  io.EOF,
	}) {
		t.Fatalf("unexpected: got=%#v", xerr)
	}
}

func TestErrorArray(t *testing.T) {
	err1 := errors.New("hello")
	err2 := errors.New("world")
	eaA := errorArray{err1, err2}
	eaB := errorArray{err2, err1}
	eaC := errorArray{parallelError{err: io.EOF}, err2}

	t.Run("Error", func(t *testing.T) {
		if s := eaA.Error(); s != "hello" {
			t.Fatalf("unexpected eaA.Error: got=%s", s)
		}
		if s := eaB.Error(); s != "world" {
			t.Fatalf("unexpected eaB.Error: got=%s", s)
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		if x := errors.Unwrap(eaA); x != nil {
			t.Fatalf("unexpected eaA.Unwrap: got=%s", x)
		}
		if x := errors.Unwrap(eaB); x != nil {
			t.Fatalf("unexpected eaB.Unwrap: got=%s", x)
		}
		if x := errors.Unwrap(eaC); x != io.EOF {
			t.Fatalf("unexpected eaC.Unwrap: got=%s", x)
		}
	})

	t.Run("Is", func(t *testing.T) {
		for i, tc := range []struct {
			ea  errorArray
			err error
			exp bool
		}{
			{eaA, err1, true},
			{eaA, err2, true},
			{eaA, io.EOF, false},
			{eaC, err1, false},
			{eaC, err2, true},
			{eaC, io.EOF, true},
		} {
			act := errors.Is(tc.ea, tc.err)
			if act != tc.exp {
				t.Fatalf("unexpected errors.Is: #%d want=%t got=%t",
					i, tc.exp, act)
			}
		}
	})

	t.Run("As", func(t *testing.T) {
		err3 := &os.PathError{Op: "dummyOp", Path: "dummyPath", Err: io.EOF}
		ea := errorArray{io.EOF, err3}
		var xerr *os.PathError
		if !errors.As(ea, &xerr) {
			t.Fatal("errors.As(os.PathError) failed")
		}
		if !reflect.DeepEqual(xerr, err3) {
			t.Fatalf("unexpected: got=%#v", xerr)
		}
		// not found case
		var xerr2 *os.PathError
		if errors.As(eaA, &xerr2) {
			t.Fatalf("unexpected errors.As succeed: %s", xerr2)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for i, tc := range []struct {
			ea  errorArray
			exp []error
		}{
			{eaA, []error{err1, err2}},
		} {
			act := tc.ea.Errors()
			if !reflect.DeepEqual(tc.exp, act) {
				t.Fatalf("unexpected Errors() #%d: want=%+v got=%+v",
					i, tc.exp, act)
			}
		}
	})
}
