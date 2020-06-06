package jobrun

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"
)

func TestSerial(t *testing.T) {
	var jobs Serial
	var stage int
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		if stage != 0 {
			t.Fatalf("unexpected stage: want=0 got=%d", stage)
		}
		stage = 1
		return nil
	}))
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		if stage != 1 {
			t.Fatalf("unexpected stage: want=1 got=%d", stage)
		}
		stage = 2
		return nil
	}))
	err := jobs.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected failure: %s", err)
	}
	if stage != 2 {
		t.Fatalf("unexpected stage: want=2 got=%d", stage)
	}
}

func TestSerial_failure(t *testing.T) {
	var jobs Serial
	var stage int
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		if stage != 0 {
			t.Fatalf("unexpected stage: want=0 got=%d", stage)
		}
		stage = 1
		return nil
	}))
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		if stage != 1 {
			t.Fatalf("unexpected stage: want=1 got=%d", stage)
		}
		stage = 2
		return io.EOF
	}))
	err := jobs.Run(context.Background())
	if err == nil {
		t.Fatal("unexpected succeed")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("unexpected failure: %s", err)
	}
	if stage != 2 {
		t.Fatalf("unexpected stage: want=2 got=%d", stage)
	}
}

func TestSerial_cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var jobs Serial
	var stage int
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		if stage != 0 {
			t.Fatalf("unexpected stage: want=0 got=%d", stage)
		}
		stage = 1
		return nil
	}))
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		if stage != 1 {
			t.Fatalf("unexpected stage: want=1 got=%d", stage)
		}
		stage = 2
		cancel()
		return nil
	}))
	err := jobs.Run(ctx)
	if err == nil {
		t.Fatal("unexpected succeed")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected failure: %s", err)
	}
	if stage != 2 {
		t.Fatalf("unexpected stage: want=2 got=%d", stage)
	}
}

func TestParallel(t *testing.T) {
	var jobs Parallel
	var state1, state2 int32
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		atomic.StoreInt32(&state1, 1)
		for atomic.LoadInt32(&state2) == 0 {
			time.Sleep(time.Millisecond)
		}
		atomic.StoreInt32(&state1, 2)
		return nil
	}))
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		atomic.StoreInt32(&state2, 1)
		for atomic.LoadInt32(&state1) == 0 {
			time.Sleep(time.Millisecond)
		}
		atomic.StoreInt32(&state2, 2)
		return nil
	}))
	err := jobs.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected failure: %s", err)
	}
	if state1 != 2 || state2 != 2 {
		t.Fatalf("jobs not complete: state1=%d state2=%d", state1, state2)
	}
}

func TestParallel_failure(t *testing.T) {
	var jobs Parallel
	var state1, state2 int32
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		atomic.StoreInt32(&state1, 1)
		for atomic.LoadInt32(&state2) == 0 {
			time.Sleep(time.Millisecond)
		}
		atomic.StoreInt32(&state1, 2)
		return nil
	}))
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		atomic.StoreInt32(&state2, 1)
		for atomic.LoadInt32(&state1) == 0 {
			time.Sleep(time.Millisecond)
		}
		atomic.StoreInt32(&state2, 2)
		return io.EOF
	}))
	err := jobs.Run(context.Background())
	if err == nil {
		t.Fatal("unexpected succeed")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("unexpected error: %s", err)
	}
	if state1 != 2 || state2 != 2 {
		t.Fatalf("jobs not complete: state1=%d state2=%d", state1, state2)
	}
}

func TestParallel_cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var jobs Parallel
	var state1, state2 int32
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		atomic.StoreInt32(&state1, 1)
		for atomic.LoadInt32(&state2) == 0 {
			time.Sleep(time.Millisecond)
		}
		atomic.StoreInt32(&state1, 2)
		ti := time.After(500 * time.Millisecond)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ti:
			atomic.StoreInt32(&state1, 3)
		}
		return nil
	}))
	jobs.Add(RunnerFunc(func(ctx context.Context) error {
		atomic.StoreInt32(&state2, 1)
		for atomic.LoadInt32(&state1) == 0 {
			time.Sleep(time.Millisecond)
		}
		atomic.StoreInt32(&state2, 2)
		cancel()
		return nil
	}))
	err := jobs.Run(ctx)
	if err == nil {
		t.Fatal("unexpected succeed")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected failure: %s", err)
	}
	if state1 != 2 || state2 != 2 {
		t.Fatalf("jobs not complete: state1=%d state2=%d", state1, state2)
	}
}
