// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/model/obs"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Job contains all state required to perform a runner operation for a given subject.
type Job struct {
	// backend is the backend used to produce the recipes being run.
	// We retain the backend to be able to work out how to parse the run results.
	backend *service.Backend

	// parser is the observation parser used to interpret the results of a run.
	parser ObsParser

	// resCh is the channel to which we're sending the run result.
	resCh chan<- builder.Request

	// subject is a pointer to the subject being run.
	subject *subject.Named

	// quantities is the set of quantities used to parametrise the running job.
	quantities QuantitySet
}

// Run runs the job with context ctx.
func (j *Job) Run(ctx context.Context) error {
	for cidstr, c := range j.subject.Compiles {
		cid := id.FromString(cidstr)
		if err := j.runCompile(ctx, cid, &c); err != nil {
			return err
		}
	}
	return nil
}

func (j *Job) runCompile(ctx context.Context, cid id.ID, c *subject.CompileResult) error {
	run, err := j.runCompileInner(ctx, cid, c)
	if err != nil {
		return err
	}
	return j.makeBuilderReq(cid, run).SendTo(ctx, j.resCh)
}

func (j *Job) runCompileInner(ctx context.Context, cid id.ID, c *subject.CompileResult) (subject.RunResult, error) {
	if !c.Status.IsOk() {
		return subject.RunResult{Result: subject.Result{Status: c.Status}}, nil
	}

	bin := c.Files.Bin
	if bin == "" {
		return subject.RunResult{
			Result: subject.Result{Status: status.Unknown},
		}, fmt.Errorf("%w: subject=%s, compiler=%s", ErrNoBin, j.subject.Name, cid.String())
	}

	start := time.Now()
	o, runErr := j.runAndParseBin(ctx, cid, bin)
	s, err := statusOfRun(o, runErr)

	return j.makeResult(start, s, o), err
}

func (j *Job) makeResult(start time.Time, s status.Status, o *obs.Obs) subject.RunResult {
	return subject.RunResult{
		Result: subject.Result{
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
func (j *Job) runAndParseBin(ctx context.Context, cid id.ID, bin string) (*obs.Obs, error) {
	tctx, cancel := j.quantities.Timeout.OnContext(ctx)
	defer cancel()

	cmd := exec.CommandContext(tctx, bin)
	obsr, err := cmd.StdoutPipe()
	if err != nil {
		return nil, j.liftError(cid, "opening pipe for", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, j.liftError(cid, "starting", err)
	}

	var o obs.Obs
	perr := j.parser.ParseObs(tctx, j.backend, obsr, &o)
	werr := cmd.Wait()

	return &o, mostRelevantError(werr, perr, tctx.Err())
}

// mostRelevantError tries to get the 'most relevant' error, given the run errors r, parsing errors p, and
// possible context errors c.
//
// The order of relevance is as follows:
// - Timeouts (through c)
// - Run errors (through r)
// - Parse errors (through p)
//
// We assume that no other context errors need to be propagated.
func mostRelevantError(r, p, c error) error {
	switch {
	case c != nil && errors.Is(c, context.DeadlineExceeded):
		return c
	case r != nil:
		return r
	default:
		return p
	}
}

func (j *Job) makeBuilderReq(cid id.ID, run subject.RunResult) builder.Request {
	return builder.RunRequest(j.subject.Name, cid, run)
}

// liftError wraps err with context about where it occurred.
func (j *Job) liftError(cid id.ID, stage string, err error) error {
	if err == nil {
		return nil
	}
	return Error{
		Stage:    stage,
		Compiler: cid,
		Subject:  j.subject.Name,
		Inner:    err,
	}
}
