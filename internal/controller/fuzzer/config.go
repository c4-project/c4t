// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"context"
	"fmt"
	"log"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// SubjectPather is the interface of things that serve file-paths for subject outputs during a fuzz batch.
type SubjectPather interface {
	// Prepare sets up the directories ready to serve through SubjectPaths.
	Prepare() error
	// SubjectLitmus gets the litmus file path for the subject/cycle pair sc.
	SubjectLitmus(sc SubjectCycle) string
	// SubjectTrace gets the trace file path for the subject/cycle pair sc.
	SubjectTrace(sc SubjectCycle) string
}

//go:generate mockery -name SubjectPather

// Config represents the configuration that goes into a batch fuzzer run.
type Config struct {
	// Driver holds the single-file fuzzer that the fuzzer is going to use.
	Driver SingleFuzzer

	// StatDumper tells the fuzzer how to scrape statistics from the fuzzed outputs.
	StatDumper litmus.StatDumper

	// Logger is the logger to use while fuzzing.  It may be nil, which silences logging.
	Logger *log.Logger

	// Observers observe the fuzzer's progress across a corpus.
	Observers []builder.Observer

	// Paths contains the path set for things generated by this fuzzer.
	Paths SubjectPather

	// Quantities sets the quantities for this batch fuzzer run.
	Quantities QuantitySet
}

// Check makes sure various parts of this config are present.
func (c *Config) Check() error {
	if c.Driver == nil {
		return ErrDriverNil
	}
	if c.StatDumper == nil {
		return ErrStatDumperNil
	}
	if c.Paths == nil {
		return iohelp.ErrPathsetNil
	}
	if c.Quantities.SubjectCycles <= 0 {
		return fmt.Errorf("%w: non-positive subject cycle amount", corpus.ErrSmall)
	}
	return nil
}

// Run runs a fuzzer configured by this config.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	f, err := New(c, p)
	if err != nil {
		return nil, err
	}
	return f.Fuzz(ctx)
}
