package jobrun

import "context"

// NamedRunner is a runner with its name.
type NamedRunner interface {
	Runner

	Name() string
}

type namedRunner struct {
	Runner

	name string
}

func (nr *namedRunner) Name() string {
	return nr.name
}

// Name names a runner.
func Name(name string, runner Runner) NamedRunner {
	return &namedRunner{
		Runner: runner,
		name:   name,
	}
}

// NameFunc names a function as runner.
func NameFunc(name string, fn func(ctx context.Context) error) NamedRunner {
	return Name(name, RunnerFunc(fn))
}

// name gets a name of runner if available. otherwise return empty string.
func name(r Runner) string {
	if nr, ok := r.(NamedRunner); ok {
		return nr.Name()
	}
	return ""
}
