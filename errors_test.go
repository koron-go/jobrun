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
	err := serialError(123, fmt.Errorf("hello world"))
	if err.Error() != "serial job #123 failed: hello world" {
		t.Fatal("unexpected Error()")
	}
}

func TestParallelErrors(t *testing.T) {
	err := parallelError{
		n:   123,
		err: fmt.Errorf("hello world"),
	}
	if err.Error() != "parallel job #123 failed: hello world" {
		t.Fatal("unexpected Error()")
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
