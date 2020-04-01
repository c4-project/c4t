// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Job represents the state of a compiler run.
type Job struct {
	// MachineID is the ID of the machine.
	MachineID id.ID

	// Compiler points to the compiler to run.
	Compiler *compiler.Named

	// Conf is the configuration with which this batch compiler was configured.
	Conf *Config

	// ResCh is the channel to which the compile run should send compiled subject records.
	ResCh chan<- builder.Request

	// Corpus is the corpus to compile.
	Corpus corpus.Corpus
}

func (j *Job) Compile(ctx context.Context) error {
	if j.Conf.Paths == nil {
		return fmt.Errorf("in job: %w", iohelp.ErrPathsetNil)
	}

	return j.Corpus.Each(func(s subject.Named) error {
		return j.compileSubject(ctx, &s)
	})
}

func (j *Job) compileSubject(ctx context.Context, s *subject.Named) error {
	h, herr := s.Harness(j.Compiler.Arch)
	if herr != nil {
		return herr
	}

	sp := j.Conf.Paths.SubjectPaths(SubjectCompile{CompilerID: j.Compiler.ID, Name: s.Name})

	res, rerr := j.runCompiler(ctx, sp, h)
	if rerr != nil {
		return rerr
	}

	return j.sendResult(ctx, s.Name, res)
}

func (j *Job) runCompiler(ctx context.Context, sp subject.CompileFileset, h subject.Harness) (subject.CompileResult, error) {
	logf, err := os.Create(sp.Log)
	if err != nil {
		return subject.CompileResult{}, err
	}

	tctx, cancel := j.timeout(ctx)
	defer cancel()

	start := time.Now()

	// Some compiler errors are recoverable, so we don't immediately bail on them.
	rerr := j.Conf.Driver.RunCompiler(tctx, j.compileJob(h, sp), logf)
	lerr := logf.Close()

	// We could close the log file here, but we want fatal compiler errors to take priority over log file close errors.
	return j.makeCompileResult(sp, start, mostRelevantError(rerr, lerr, tctx.Err()))
}

func (j *Job) compileJob(h subject.Harness, sp subject.CompileFileset) job.Compile {
	return job.Compile{
		In:       h.CPaths(),
		Out:      sp.Bin,
		Compiler: &j.Compiler.Compiler,
	}
}

// makeCompileResult makes a compile result given a possible err and fileset sp.
// It fails if the error is considered substantially fatal.
func (j *Job) makeCompileResult(sp subject.CompileFileset, start time.Time, err error) (subject.CompileResult, error) {
	cr := subject.CompileResult{
		Result: subject.Result{
			Time:     start,
			Duration: time.Since(start),
			Status:   subject.StatusUnknown,
		},
		Files: sp.StripMissing(),
	}

	cr.Status, err = subject.StatusOfCompileError(err)
	return cr, err
}

// sendResult tries to send a compile job result to the result channel.
// If the context ctx has been cancelled, it will fail and instead terminate the job.
func (j *Job) sendResult(ctx context.Context, name string, r subject.CompileResult) error {
	return builder.CompileRequest(name, j.Compiler.ID, r).SendTo(ctx, j.ResCh)
}

func (j *Job) timeout(ctx context.Context) (context.Context, context.CancelFunc) {
	// TODO(@MattWindsor91): dedupe with runner equivalent
	if j.Conf.Timeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, j.Conf.Timeout)
}

// mostRelevantError tries to get the 'most relevant' error, given the run errors r, parsing errors p, and
// possible context errors c.
//
// The order of relevance is as follows:
// - Timeouts (through c)
// - Run errors (through r)
// - Log file close errors (through l)
//
// We assume that no other context errors need to be propagated.
func mostRelevantError(r, l, c error) error {
	// TODO(@MattWindsor91): dedupe with runner equivalent
	switch {
	case c != nil && errors.Is(c, context.DeadlineExceeded):
		return c
	case r != nil:
		return r
	default:
		return l
	}
}
