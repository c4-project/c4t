// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/quantity"

	"github.com/MattWindsor91/act-tester/internal/model/job/compile"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Instance represents the state of a single per-subject instance of the batch compiler.
type Instance struct {
	// machineID is the ID of the machine.
	machineID id.ID

	// subject is the subject to compile.
	subject subject.Named

	// compilers points to the compilers to run.
	compilers map[string]compiler.Compiler

	// driver tells the instance how to run the compiler.
	driver Driver

	// paths tells the instance which paths to use.
	paths SubjectPather

	// quantities is the quantity set for this instance.
	quantities quantity.SingleSet

	// resCh is the channel to which the compile run should send compiled subject records.
	resCh chan<- builder.Request
}

func (j *Instance) Compile(ctx context.Context) error {
	if j.paths == nil {
		return fmt.Errorf("in job: %w", iohelp.ErrPathsetNil)
	}

	for n, c := range j.compilers {
		nc, err := c.AddNameString(n)
		if err != nil {
			return err
		}
		if err := j.compileOnCompiler(ctx, nc); err != nil {
			return err
		}
	}
	return nil
}

func (j *Instance) compileOnCompiler(ctx context.Context, nc *compiler.Named) error {
	h, herr := j.subject.Recipe(nc.Arch)
	if herr != nil {
		return herr
	}

	sc := compilation.Name{CompilerID: nc.ID, SubjectName: j.subject.Name}
	sp := j.paths.SubjectPaths(sc)

	res, rerr := j.runCompiler(ctx, nc, sp, h)
	if rerr != nil {
		return rerr
	}

	return builder.CompileRequest(sc, res).SendTo(ctx, j.resCh)
}

func (j *Instance) runCompiler(ctx context.Context, nc *compiler.Named, sp compilation.CompileFileset, h recipe.Recipe) (compilation.CompileResult, error) {
	logf, err := j.openLogFile(sp.Log)
	if err != nil {
		return compilation.CompileResult{}, err
	}

	tctx, cancel := j.quantities.Timeout.OnContext(ctx)
	defer cancel()

	start := time.Now()

	job := j.compileJob(h, nc, sp)
	// Some compiler errors are recoverable, so we don't immediately bail on them.
	rerr := j.runCompilerJob(tctx, job, logf)

	lerr := logf.Close()

	// We could close the log file here, but we want fatal compiler errors to take priority over log file close errors.
	return j.makeCompileResult(sp, start, mostRelevantError(rerr, lerr, tctx.Err()))
}

func (j *Instance) runCompilerJob(ctx context.Context, job compile.Recipe, logf io.Writer) error {
	i, err := NewInterpreter(j.driver, job, ILogTo(logf))
	if err != nil {
		return err
	}
	return i.Interpret(ctx)
}

func (j *Instance) openLogFile(l string) (io.WriteCloser, error) {
	if ystring.IsBlank(l) {
		return iohelp.DiscardCloser(), nil
	}
	return os.Create(l)
}

func (j *Instance) compileJob(r recipe.Recipe, nc *compiler.Named, sp compilation.CompileFileset) compile.Recipe {
	return compile.FromRecipe(&nc.Compiler, r, sp.Bin)
}

// makeCompileResult makes a compile result given a possible err and fileset sp.
// It fails if the error is considered substantially fatal.
func (j *Instance) makeCompileResult(sp compilation.CompileFileset, start time.Time, err error) (compilation.CompileResult, error) {
	cr := compilation.CompileResult{
		Result: compilation.Result{
			Time:     start,
			Duration: time.Since(start),
			Status:   status.Unknown,
		},
		Files: sp.StripMissing(),
	}

	cr.Status, err = status.FromCompileError(err)
	return cr, err
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
