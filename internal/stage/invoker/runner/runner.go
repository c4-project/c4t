// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package runner contains low-level code for running the machine node via SSH and locally.
package runner

import (
	"context"

	"github.com/MattWindsor91/c4t/internal/quantity"

	"github.com/MattWindsor91/c4t/internal/plan"
	"github.com/MattWindsor91/c4t/internal/remote"
)

// Runner is the interface of types that know how to run the machine node.
type Runner interface {
	// Send performs any copying and transformation needed for p to run.
	// It returns a pointer to the plan to send to the machine node, which may or may not be p.
	Send(ctx context.Context, p *plan.Plan) (*plan.Plan, error)

	// Start starts the machine binary, returning a set of pipe readers and writers to use for communication with it.
	Start(ctx context.Context, qs quantity.MachNodeSet) (*remote.Pipeset, error)

	// Wait blocks waiting for the command to finish (or the context passed into Start to cancel).
	Wait() error

	// Recv merges the post-run plan runp into the original plan origp, copying back any files needed.
	// It returns a pointer to the final 'merged', which may or may not be origp and runp.
	// It may modify origp in place.
	Recv(ctx context.Context, origp, runp *plan.Plan) (*plan.Plan, error)
}
