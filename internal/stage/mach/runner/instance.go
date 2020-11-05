// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	backend3 "github.com/MattWindsor91/act-tester/internal/model/service/backend"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"

	"github.com/MattWindsor91/act-tester/internal/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/subject/status"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/subject/obs"

	"github.com/MattWindsor91/act-tester/internal/subject"
)

// Instance contains all state required to perform a runner operation for a given subject.
type Instance struct {
	// backend is the backend used to produce the recipes being run.
	// We retain the backend to be able to work out how to parse the run results.
	backend *backend3.Spec

	// parser is the observation parser used to interpret the results of a run.
	parser ObsParser

	// resCh is the channel to which we're sending the run result.
	resCh chan<- builder.Request

	// subject is a pointer to the subject being run.
	subject *subject.Named

	// quantities is the set of quantities used to parametrise the running job.
	quantities quantity.BatchSet
}

// Run runs the instance with context ctx.
func (n *Instance) Run(ctx context.Context) error {
	for cidstr, cc := range n.subject.Compilations {
		cid := id.FromString(cidstr)
		name := compilation.Name{CompilerID: cid, SubjectName: n.subject.Name}
		if err := n.runCompile(ctx, name, cc.Compile); err != nil {
			return err
		}
	}
	return nil
}

func (n *Instance) runCompile(ctx context.Context, name compilation.Name, c *compilation.CompileResult) error {
	if c == nil {
		return fmt.Errorf("%w: %s", subject.ErrMissingCompile, name)
	}
	run, err := n.runCompileInner(ctx, name, c)
	if err != nil {
		return err
	}
	return builder.RunRequest(name, run).SendTo(ctx, n.resCh)
}

func (n *Instance) runCompileInner(ctx context.Context, name compilation.Name, c *compilation.CompileResult) (compilation.RunResult, error) {
	if !c.Status.IsOk() {
		return compilation.RunResult{Result: compilation.Result{Status: c.Status}}, nil
	}

	bin := c.Files.Bin
	if bin == "" {
		return compilation.RunResult{
			Result: compilation.Result{Status: status.Unknown},
		}, fmt.Errorf("%w: %s", ErrNoBin, name)
	}

	start := time.Now()
	o, runErr := n.runAndParseBin(ctx, name, bin)
	s, err := statusOfRun(o, runErr)

	return n.makeResult(start, s, o), err
}

func (n *Instance) makeResult(start time.Time, s status.Status, o *obs.Obs) compilation.RunResult {
	return compilation.RunResult{
		Result: compilation.Result{
			Time:     start,
			Duration: time.Since(start),
			Status:   s,
		},
		Obs: o,
	}
}

func statusOfRun(o *obs.Obs, runErr error) (status.Status, error) {
	if runErr != nil {
		return status.FromRunError(runErr)
	}
	return o.Status(), nil
}

// runAndParseBin runs the binary at bin and parses its result into an observation struct.
func (n *Instance) runAndParseBin(ctx context.Context, name compilation.Name, bin string) (*obs.Obs, error) {
	tctx, cancel := n.quantities.Timeout.OnContext(ctx)
	defer cancel()

	cmd := exec.CommandContext(tctx, bin)
	obsr, err := cmd.StdoutPipe()
	if err != nil {
		return nil, n.liftError(name, "opening pipe for", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, n.liftError(name, "starting", err)
	}

	var o obs.Obs
	perr := n.parser.ParseObs(tctx, n.backend, obsr, &o)
	werr := cmd.Wait()

	return &o, errhelp.TimeoutOrFirstError(tctx, werr, perr)
}

// liftError wraps err with context about where it occurred.
func (n *Instance) liftError(name compilation.Name, stage string, err error) error {
	if err == nil {
		return nil
	}
	return Error{
		Stage:       stage,
		Compilation: name,
		Inner:       err,
	}
}
