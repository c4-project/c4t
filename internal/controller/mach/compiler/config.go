// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// SingleRunner is the interface of things that can run compilers.
type SingleRunner interface {
	// RunCompiler runs the compiler job j.
	// If applicable, errw will be connected to the compiler's standard error.
	RunCompiler(ctx context.Context, j job.Compile, errw io.Writer) error
}

// SubjectPather is the interface of types that can produce path sets for compilations.
type SubjectPather interface {
	// Prepare sets up the directories ready to serve through SubjectPaths.
	// It takes the list of compiler IDs that are to be represented in the pathset.
	Prepare(compilers []id.ID) error

	// SubjectPaths gets the binary and log file paths for the subject/compiler pair sc.
	SubjectPaths(sc SubjectCompile) subject.CompileFileset
}

// Config represents the configuration that goes into a batch compiler run.
type Config struct {
	// Timeout is the timeout for each compile.
	// Non-positive values disable the timeout.
	Timeout time.Duration

	// Driver is what the compiler should use to run single compiler jobs.
	Driver SingleRunner

	// Observers observe the compiler's progress across a corpus.
	Observers []builder.Observer

	// Logger is the logger used for informational output during the compile.
	Logger *log.Logger

	// Paths is the pathset for this compiler run.
	Paths SubjectPather
}

// Check checks for various problems with a config.
func (c *Config) Check() error {
	if c.Driver == nil {
		return ErrDriverNil
	}
	if c.Paths == nil {
		return iohelp.ErrPathsetNil
	}
	return nil
}

// Run runs a compiler configured by this config.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	cm, err := New(c, p)
	if err != nil {
		return nil, err
	}
	return cm.Run(ctx)
}
