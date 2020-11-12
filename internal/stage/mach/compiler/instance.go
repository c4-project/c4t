// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/interpreter"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"

	"github.com/MattWindsor91/act-tester/internal/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/subject/status"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/subject"
)

// Instance represents the state of a single per-subject instance of the batch compiler.
type Instance struct {
	// machineID is the ID of the machine.
	machineID id.ID

	// subject is the subject to compile.
	subject subject.Named

	// compilers points to the compilers to run.
	compilers map[string]compiler.Configuration

	// driver tells the instance how to run the compiler.
	driver interpreter.Driver

	// paths tells the instance which paths to use.
	paths SubjectPather

	// quantities is the quantity set for this instance.
	quantities quantity.BatchSet

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
	rid, r, err := j.subject.Recipe(nc.Arch)
	if err != nil {
		return err
	}

	sc := compilation.Name{CompilerID: nc.ID, SubjectName: j.subject.Name}
	res := compilation.CompileResult{
		Result: compilation.Result{
			Status: status.Unknown,
		},
		RecipeID: rid,
		Files:    j.paths.SubjectPaths(sc),
	}

	if r.NeedsCompile() {
		if err := j.runCompiler(ctx, nc, &res, r); err != nil {
			return err
		}
	}
	res.Files = res.Files.StripMissing()

	return builder.CompileRequest(sc, res).SendTo(ctx, j.resCh)
}

func (j *Instance) runCompiler(ctx context.Context, nc *compiler.Named, res *compilation.CompileResult, h recipe.Recipe) error {
	logf, err := j.openLogFile(res.Files.Log)
	if err != nil {
		return err
	}

	tctx, cancel := j.quantities.Timeout.OnContext(ctx)
	defer cancel()

	res.Time = time.Now()

	// Some compiler errors are recoverable, so we don't immediately bail on them.
	rerr := j.runCompilerJob(tctx, nc, res.Files, h, logf)
	lerr := logf.Close()

	res.Duration = time.Since(res.Time)
	res.Status, err = status.FromCompileError(errhelp.TimeoutOrFirstError(tctx, rerr, lerr))
	return err
}

func (j *Instance) runCompilerJob(ctx context.Context, nc *compiler.Named, sp compilation.CompileFileset, h recipe.Recipe, logf io.Writer) error {

	i, err := interpreter.NewInterpreter(j.driver, &nc.Configuration, sp.Bin, h, interpreter.LogTo(logf))
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
