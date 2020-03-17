package jobrun

import (
	"fmt"
	"testing"
)

func TestParallelErrors(t *testing.T) {
	err := parallelError{
		n:   123,
		err: fmt.Errorf("hello world"),
	}
	if err.Error() != "parallel job #123 failed: hello world" {
		t.Fatal("unexpected Error()")
	}
}
