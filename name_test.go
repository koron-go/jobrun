package jobrun

import (
	"context"
	"errors"
	"testing"
)

func TestNameSerial(t *testing.T) {
	var stage int
	jobs := Serial{
		NameFunc("foo", func(ctx context.Context) error {
			if stage != 0 {
				t.Fatalf("unexpected stage: want=0 got=%d", stage)
			}
			stage = 1
			return nil
		}),
		NameFunc("bar", func(ctx context.Context) error {
			if stage != 1 {
				t.Fatalf("unexpected stage: want=1 got=%d", stage)
			}
			stage = 2
			return errors.New("expected failure")
		}),
		NameFunc("quux", func(ctx context.Context) error {
			t.Fatal("never called")
			return nil
		}),
	}
	err := jobs.Run(context.Background())
	if err == nil {
		t.Fatal("unexpected succeed")
	}
	if stage != 2 {
		t.Errorf("unexpected stage: want=2 got=%d", stage)
	}
	if got := err.Error(); got != "job:bar (serial job #1) failed: expected failure" {
		t.Errorf("unexpected error details: got=%s", got)
	}
}

func TestNameParallel(t *testing.T) {
	jobs := Parallel{
		NameFunc("foo", func(ctx context.Context) error {
			return nil
		}),
		NameFunc("bar", func(ctx context.Context) error {
			return errors.New("expected failure")
		}),
		NameFunc("quux", func(ctx context.Context) error {
			return nil
		}),
	}
	err := jobs.Run(context.Background())
	if err == nil {
		t.Fatal("unexpected succeed")
	}
	if got := err.Error(); got != "job:bar (parallel job #1) failed: expected failure" {
		t.Errorf("unexpected error details: got=%s", got)
	}
}
