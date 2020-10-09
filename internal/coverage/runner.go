// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/MattWindsor91/act-tester/internal/subject"
)

var (
	// ErrNeedRunInfo occurs when we try to instantiate a runner for a standalone profile without run information.
	ErrNeedRunInfo = errors.New("need run information for this profile")

	// ErrUnsupportedProfileKind occurs when we try to instantiate a runner for an unsupported profile type.
	ErrUnsupportedProfileKind = errors.New("this profile kind can't be run yet")
)

// RunnerContext is the type of state provided to a coverage runner.
type RunnerContext struct {
	// Seed is the seed to use to drive any random parts of the coverage runner.
	Seed int32
	// BucketDir is the filepath to the bucket directory into which the coverage runner should output its recipe.
	BucketDir string
	// NumInBucket is the index of this single instance in its bucket.
	NumInBucket int
	// Input points to an input subject for the coverage runner, if any are available.
	Input *subject.Subject
}

// inputPath tries to get the filepath to the currently available input's litmus test.
// It returns the empty string if no such file is available.
func (r RunnerContext) inputPath() string {
	if r.Input == nil {
		return ""
	}
	l, err := r.Input.BestLitmus()
	if err != nil {
		return ""
	}
	return filepath.Clean(l.Path)
}

// outputPath gets the filepath to which the runner should output one C file.
func (r RunnerContext) outputPath() string {
	return filepath.Join(r.BucketDir, fmt.Sprintf("%d.c", r.NumInBucket))
}

// ExpandArgs expands various special identifiers in args to parts of the runner context.
func (r RunnerContext) ExpandArgs(arg ...string) []string {
	nargs := make([]string, len(arg))
	for i, a := range arg {
		nargs[i] = r.expandArg(a)

	}
	return nargs
}

func (r RunnerContext) expandArg(arg string) string {
	switch arg {
	case "$seed":
		return strconv.Itoa(int(r.Seed))
	case "$input":
		return r.inputPath()
	case "$output":
		return r.outputPath()
	default:
		return arg
	}
}

// Runner is the interface of things that can be run to generate coverage testbeds.
type Runner interface {
	// Run runs the Runner with context ctx and runner context rc.
	Run(ctx context.Context, rc RunnerContext) error
}

//go:generate mockery --name=Runner

// StandaloneRunner is a coverage runner that runs a standalone binary.
type StandaloneRunner struct {
	// run tells the runner how to run the standalone runner.
	run service.RunInfo
	// errw is the writer to which stderr should go, if any.
	errw io.Writer
}

// Run runs the standalone runner.
func (s *StandaloneRunner) Run(ctx context.Context, rc RunnerContext) error {
	cmd := exec.CommandContext(ctx, s.run.Cmd, rc.ExpandArgs(s.run.Args...)...)
	cmd.Stderr = s.errw
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running coverage generator %q: %w", s.run.Cmd, err)
	}
	return nil
}

func (p *Profile) runner(errw io.Writer) (Runner, error) {
	// this mostly used only for testing
	if p.Runner != nil {
		return p.Runner, nil
	}

	switch p.Kind {
	case Standalone:
		return p.standaloneRunner(errw)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProfileKind, p.Kind)
	}
}

func (p *Profile) standaloneRunner(errw io.Writer) (*StandaloneRunner, error) {
	if p.Run == nil {
		return nil, ErrNeedRunInfo
	}
	return &StandaloneRunner{run: *p.Run, errw: errw}, nil
}
