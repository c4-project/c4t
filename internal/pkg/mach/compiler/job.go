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
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/job"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/service"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/helpers/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"
)

// Job represents the state of a compiler run.
type Job struct {
	// MachineID is the ID of the machine.
	MachineID id.ID

	// Compiler points to the compiler to run.
	Compiler *service.NamedCompiler

	// Pathset is the pathset to use for this compiler run.
	Pathset SubjectPather

	// Compiler is the compiler runner to use for running compilers.
	Runner SingleRunner

	// ResCh is the channel to which the compile run should send compiled subject records.
	ResCh chan<- builder.Request

	// Corpus is the corpus to compile.
	Corpus corpus.Corpus
}

func (j *Job) Compile(ctx context.Context) error {
	if j.Pathset == nil {
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

	sp := j.Pathset.SubjectPaths(SubjectCompile{CompilerID: j.Compiler.ID, Name: s.Name})

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

	// Some compiler errors are recoverable, so we don't immediately bail on them.
	cerr := j.Runner.RunCompiler(ctx, j.compileJob(h, sp), logf)

	// We could close the log file here, but we want fatal compiler errors to take priority over log file close errors.
	res, rerr := j.makeCompileResult(sp, cerr)
	if rerr != nil {
		_ = logf.Close()
		return subject.CompileResult{}, rerr
	}

	lerr := logf.Close()
	return res, lerr
}

func (j *Job) compileJob(h subject.Harness, sp subject.CompileFileset) job.Compile {
	return job.Compile{
		In:       h.CPaths(),
		Out:      sp.Bin,
		Compiler: &j.Compiler.Compiler,
	}
}

// makeCompileResult makes a compile result given a possible compile error cerr and fileset sp.
// It fails if the compile error is considered substantially fatal.
func (j *Job) makeCompileResult(sp subject.CompileFileset, cerr error) (subject.CompileResult, error) {
	cr := subject.CompileResult{
		// Potentially overridden further down.
		Success: false,
		Files:   sp.StripMissing(),
	}

	if cerr != nil {
		// If the error was the compiler process failing to run, then we should report that and carry on.
		var perr *exec.ExitError
		if !errors.As(cerr, &perr) {
			return cr, cerr
		}

		return cr, nil
	}

	cr.Success = true
	return cr, nil
}

// sendResult tries to send a compile job result to the result channel.
// If the context ctx has been cancelled, it will fail and instead terminate the job.
func (j *Job) sendResult(ctx context.Context, name string, r subject.CompileResult) error {
	return builder.CompileRequest(name, j.Compiler.ID, r).SendTo(ctx, j.ResCh)
}
