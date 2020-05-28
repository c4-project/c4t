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

	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/forward"
	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// RunnerFactory is the interface of sources of machine node runners.
//
// Runner factories can contain disposable state (for example, long-running SSH connections), and so can be closed.
type RunnerFactory interface {
	// MakeRunner creates a new Runner, representing a particular invoker session on a machine.
	// It takes the plan in case the factory is waiting to get machine configuration from it.
	MakeRunner(p *plan.Plan, obs ...remote.CopyObserver) (Runner, error)

	// Runner spawners can be closed once no more runners are needed.
	// For SSH runner spawners, this will close the SSH connection.
	io.Closer
}

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
func (m *Invoker) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	runner, err := m.rfac.MakeRunner(p, m.observers.Copy...)
	if err != nil {
		return nil, fmt.Errorf("while spawning runner: %w", err)
	}

	rp, err := runner.Send(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("while copying files to machine: %w", err)
	}

	ps, err := runner.Start(ctx, m.invoker)
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

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return p.Check()
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
		Observers: m.observers.Corpus,
	}
	return rp.Run(ctx)
}

// sendPlan sends p to w, then closes w, reporting any relevant errors.
func sendPlan(p *plan.Plan, w io.WriteCloser) error {
	terr := p.Write(w)
	ierr := w.Close()
	if terr != nil {
		return fmt.Errorf("while sending input plan: %w", terr)
	}
	if ierr != nil {
		return fmt.Errorf("while closing input pipe: %w", ierr)
	}
	return nil
}
