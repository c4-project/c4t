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

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/obs"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"
)

// Job contains all state required to perform a runner operation for a given subject.
type Job struct {
	// MachConfig points to the runner config.
	Conf *Config

	// Backend is the backend used to produce the harnesses being run.
	Backend *model.Backend

	// ResCh is the channel to which we're sending the run result.
	ResCh chan<- builder.Request

	// Subject is a pointer to the subject being run.
	Subject *subject.Named
}

// Run runs the job with context ctx.
func (j *Job) Run(ctx context.Context) error {
	for cidstr, c := range j.Subject.Compiles {
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
	return j.makeBuilderReq(cid, run).SendTo(ctx, j.ResCh)
}

func (j *Job) runCompileInner(ctx context.Context, cid id.ID, c *subject.CompileResult) (subject.Run, error) {
	if !c.Success {
		return subject.Run{Status: subject.StatusCompileFail}, nil
	}

	bin := c.Files.Bin
	if bin == "" {
		return subject.Run{Status: subject.StatusUnknown}, fmt.Errorf("%w: subject=%s, compiler=%s", ErrNoBin, j.Subject.Name, cid.String())
	}

	o, runErr := j.runAndParseBin(ctx, cid, bin)
	status, err := subject.StatusOfObs(o, runErr)

	return subject.Run{Status: status, Obs: o}, err
}

// runAndParseBin runs the binary at bin and parses its result into an observation struct.
func (j *Job) runAndParseBin(ctx context.Context, cid id.ID, bin string) (*obs.Obs, error) {
	// TODO(@MattWindsor91): make the timeout configurable
	tctx, cancel := j.timeout(ctx)
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
	perr := j.Conf.Parser.ParseObs(tctx, *j.Backend, obsr, &o)
	werr := cmd.Wait()

	return &o, mostRelevantError(werr, perr, tctx.Err())
}

func (j *Job) timeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if j.Conf.Timeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, time.Duration(j.Conf.Timeout)*time.Minute)
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

func (j *Job) makeBuilderReq(cid id.ID, run subject.Run) builder.Request {
	return builder.RunRequest(j.Subject.Name, cid, run)
}

// liftError wraps err with context about where it occurred.
func (j *Job) liftError(cid id.ID, stage string, err error) error {
	if err == nil {
		return nil
	}
	return Error{
		Stage:    stage,
		Compiler: cid,
		Subject:  j.Subject.Name,
		Inner:    err,
	}
}
