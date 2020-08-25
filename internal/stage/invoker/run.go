// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/forward"
	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Run runs the machine invoker stage.
func (m *Invoker) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}
	return p.RunStage(ctx, stage.Invoke, m.invoke)
}

// invoke runs the machine binary.
func (m *Invoker) invoke(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	runner, err := m.rfac.MakeRunner(m.ldir, p, m.copyObservers...)
	if err != nil {
		return nil, fmt.Errorf("while spawning runner: %w", err)
	}

	rp, err := runner.Send(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("while copying files to machine: %w", err)
	}

	qs := m.baseQuantities
	if err := m.pqo.OverrideQuantitiesFromPlan(p, &qs); err != nil {
		return nil, fmt.Errorf("while extracting quantities from plan: %w", err)
	}
	ps, err := runner.Start(ctx, qs)
	if err != nil {
		return nil, fmt.Errorf("while starting command: %w", err)
	}

	np, err := m.runPipework(ctx, rp, ps)
	// Waiting _should_ close the pipes.
	werr := runner.Wait()

	if err != nil {
		return nil, err
	}
	if werr != nil {
		return nil, werr
	}

	return runner.Recv(ctx, p, np)
}

// Close closes any persistent connections used by this invoker.
func (m *Invoker) Close() error {
	return m.rfac.Close()
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	if err := p.Check(); err != nil {
		return err
	}
	return p.Metadata.RequireStage(stage.Plan, stage.Lift)
}

// runPipework runs the various parallel processes that read to and write from the machine binary via ps.
// These include: sending the remote plan rp to stdin; receiving the updated plan from stdout; and replaying
// observations from stderr.
func (m *Invoker) runPipework(ctx context.Context, rp *plan.Plan, ps *remote.Pipeset) (*plan.Plan, error) {
	var p2 plan.Plan

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return sendPlan(rp, ps.Stdin)
	})
	eg.Go(func() error {
		if err := plan.Read(ps.Stdout, &p2); err != nil {
			return fmt.Errorf("while decoding the output plan: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		return m.runReplayer(ectx, ps.Stderr)
	})

	return &p2, eg.Wait()
}

// runReplayer constructs and runs an observation replayer on top of r.
func (m *Invoker) runReplayer(ctx context.Context, r io.Reader) error {
	rp := forward.Replayer{
		Decoder:   json.NewDecoder(r),
		Observers: m.machObservers,
	}
	return rp.Run(ctx)
}

// sendPlan sends p to w, then closes w, reporting any relevant errors.
func sendPlan(p *plan.Plan, w io.WriteCloser) error {
	terr := p.Write(w, plan.WriteNone) // for now
	ierr := w.Close()
	if terr != nil {
		return fmt.Errorf("while sending input plan: %w", terr)
	}
	if ierr != nil {
		return fmt.Errorf("while closing input pipe: %w", ierr)
	}
	return nil
}
