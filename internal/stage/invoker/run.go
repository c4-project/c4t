// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/c4-project/c4t/internal/helper/errhelp"
	"github.com/c4-project/c4t/internal/stage/invoker/runner"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/plan/stage"

	"github.com/c4-project/c4t/internal/remote"

	"github.com/c4-project/c4t/internal/stage/mach/forward"
	"golang.org/x/sync/errgroup"

	"github.com/c4-project/c4t/internal/plan"
)

// Run runs the machine invoker stage.
func (m *Invoker) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := m.checkPlan(p); err != nil {
		return nil, err
	}
	return p.RunStage(ctx, stage.Invoke, m.invoke)
}

// invoke runs the machine binary.
func (m *Invoker) invoke(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	run, err := m.rfac.MakeRunner(m.ldir, p, m.copyObservers...)
	if err != nil {
		return nil, fmt.Errorf("while spawning runner: %w", err)
	}
	rp, err := run.Send(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("while copying files to machine: %w", err)
	}
	qs, err := m.calcQuantities(p)
	if err != nil {
		return nil, err
	}
	ps, err := run.Start(ctx, qs)
	if err != nil {
		return nil, fmt.Errorf("while starting command: %w", err)
	}
	np, err := m.awaitResults(ctx, rp, ps, run)
	if err != nil {
		return nil, err
	}
	return run.Recv(ctx, p, np)
}

func (m *Invoker) awaitResults(ctx context.Context, rp *plan.Plan, ps *remote.Pipeset, runner runner.Runner) (*plan.Plan, error) {
	np, err := m.runPipework(ctx, rp, ps)
	// Waiting _should_ close the pipes.
	werr := runner.Wait()
	return np, errhelp.FirstError(err, werr)
}

func (m *Invoker) calcQuantities(p *plan.Plan) (quantity.MachNodeSet, error) {
	qs := m.baseQuantities
	if err := m.pqo.OverrideQuantitiesFromPlan(p, &qs); err != nil {
		return qs, fmt.Errorf("while extracting quantities from plan: %w", err)
	}
	return qs, nil
}

// Close closes any persistent connections used by this invoker.
func (m *Invoker) Close() error {
	return m.rfac.Close()
}

func (m *Invoker) checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	if err := p.Check(); err != nil {
		return err
	}
	if err := p.Metadata.RequireStage(stage.Plan, stage.Lift); err != nil {
		return err
	}
	return m.handlePossibleReinvoke(p)
}

func (m *Invoker) handlePossibleReinvoke(p *plan.Plan) error {
	err := p.Metadata.ForbidStage(stage.Invoke)
	if err == nil {
		return nil
	}
	if m.allowReinvoke {
		// TODO(@MattWindsor91): strip previous invoke/compile/run metadata?
		p.Corpus.EraseCompilations()
		return nil
	}
	return err
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
