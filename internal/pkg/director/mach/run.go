// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/MattWindsor91/act-tester/internal/pkg/mach/forward"
	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

// Runner is the interface that the local and SSH runners have in common.
type Runner interface {
	// Start starts the machine binary, returning a set of pipe readers and writers to use for communication with it.
	Start(ctx context.Context) (*Pipeset, error)

	// Wait blocks waiting for the command to finish.
	Wait() error
}

// Run runs the machine binary on p.
// It presumes that p has already been amended
func (m *Mach) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	eg, ectx := errgroup.WithContext(ctx)

	ps, err := m.runner.Start(ectx)
	if err != nil {
		return nil, fmt.Errorf("while starting command: %w", err)
	}

	var p2 plan.Plan
	eg.Go(func() error {
		return sendPlan(p, ps.Stdin)
	})
	eg.Go(func() error {
		if _, err := toml.DecodeReader(ps.Stdout, &p2); err != nil {
			return fmt.Errorf("while decoding the output plan: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		r := forward.Replayer{
			Decoder:  json.NewDecoder(ps.Stderr),
			Observer: m.observer,
		}
		return r.Run(ectx)
	})

	// Waiting _should_ close the pipes.
	err = eg.Wait()
	werr := m.runner.Wait()

	if err != nil {
		return nil, err
	}
	return &p2, werr
}

// binName is the name of the machine-runner binary.
const binName = "act-tester-mach"

// runArgs produces the arguments for an invocation of binName, given directory dir.
func runArgs(dir string) []string {
	return []string{
		"-J",      // use JSON
		"-d", dir, // output to the given directory
	}
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
