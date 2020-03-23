// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"context"
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/helpers/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// SingleRunner is the interface of things that can run compilers.
type SingleRunner interface {
	// RunCompiler runs the compiler pointed to by c on the input files in j.In.
	// On success, it outputs a binary to j.Out.
	// If applicable, errw will be connected to the compiler's standard error.
	RunCompiler(ctx context.Context, c *model.NamedCompiler, j model.CompileJob, errw io.Writer) error
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
	// Driver is what the compiler should use to run single compiler jobs.
	Driver SingleRunner

	// Observer observes the compiler's progress across a corpus.
	Observer builder.Observer

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
