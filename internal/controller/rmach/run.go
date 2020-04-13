// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/transfer/remote"

	"github.com/BurntSushi/toml"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/forward"
	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Runner is the interface that the local and SSH runners have in common.
type Runner interface {
	// Send performs any copying and transformation needed for p to run.
	// It returns a pointer to the plan to send to the machine runner, which may or may not be p.
	Send(ctx context.Context, p *plan.Plan) (*plan.Plan, error)

	// Start starts the machine binary, returning a set of pipe readers and writers to use for communication with it.
	Start(ctx context.Context, mi InvocationGetter) (*remote.Pipeset, error)

	// Wait blocks waiting for the command to finish (or the context passed into Start to cancel).
	Wait() error

	// Recv merges the post-run plan runp into the original plan origp, copying back any files needed.
	// It returns a pointer to the final 'merged', which may or may not be origp and runp.
	// It may modify origp in place.
	Recv(ctx context.Context, origp, runp *plan.Plan) (*plan.Plan, error)
}

// Run runs the machine binary.
func (m *RMach) Run(ctx context.Context) (*plan.Plan, error) {
	rp, err := m.runner.Send(ctx, m.plan)
	if err != nil {
		return nil, fmt.Errorf("while copying files to machine: %w", err)
	}

	ps, err := m.runner.Start(ctx, m.conf.Invoker)
	if err != nil {
		return nil, fmt.Errorf("while starting command: %w", err)
	}

	np, err := m.runPipework(ctx, rp, ps)
	werr := m.runner.Wait()

	if err != nil {
		return nil, err
	}
	if werr != nil {
		return nil, werr
	}

	return m.runner.Recv(ctx, m.plan, np)
}

// runPipework runs the various parallel processes that read to and write from the machine binary via ps.
// These include: sending the remote plan rp to stdin; receiving the updated plan from stdout; and replaying
// observations from stderr.
func (m *RMach) runPipework(ctx context.Context, rp *plan.Plan, ps *remote.Pipeset) (*plan.Plan, error) {
	var p2 plan.Plan

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return sendPlan(rp, ps.Stdin)
	})
	eg.Go(func() error {
		if _, err := toml.DecodeReader(ps.Stdout, &p2); err != nil {
			return fmt.Errorf("while decoding the output plan: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		return m.runReplayer(ectx, ps.Stderr)
	})

	// Waiting _should_ close the pipes.
	return &p2, eg.Wait()
}

// runReplayer constructs and runs an observation replayer on top of r.
func (m *RMach) runReplayer(ctx context.Context, r io.Reader) error {
	rp := forward.Replayer{
		Decoder:   json.NewDecoder(r),
		Observers: m.conf.Observers.Corpus,
	}
	return rp.Run(ctx)
}

// sendPlan sends p to w, then closes w, reporting any relevant errors.
func sendPlan(p *plan.Plan, w io.WriteCloser) error {
	terr := p.Dump(w)
	ierr := w.Close()
	if terr != nil {
		return fmt.Errorf("while sending input plan: %w", terr)
	}
	if ierr != nil {
		return fmt.Errorf("while closing input pipe: %w", ierr)
	}
	return nil
}
