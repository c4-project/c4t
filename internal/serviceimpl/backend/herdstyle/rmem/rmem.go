// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package rmem implements rudimentary backend support for RMEM.
//
// Presently, rmem is implemented as a herdtools-style backend, despite not being a herdtools project.
// This will likely change later on.
package rmem

import (
	"context"
	"fmt"
	"io"

	"github.com/c4-project/c4t/internal/serviceimpl/backend"

	"github.com/c4-project/c4t/internal/model/service"
	backend2 "github.com/c4-project/c4t/internal/model/service/backend"
)

var armArgs = [...]string{
	"-model", "promising",
	"-model", "promise_first",
	"-model", "promising_parallel_thread_state_search",
	"-model", "promising_parallel_without_follow_trace",
	"-priority_reduction", "false",
	"-interactive", "false",
	"-hash_prune", "false",
	"-allow_partial", "true",
	"-loop_limit", "2",
}

// Rmem holds implementations of various backend responsiblities for Rmem.
type Rmem struct{}

func (Rmem) LiftStandalone(ctx context.Context, j backend2.LiftJob, r service.RunInfo, x service.Runner, w io.Writer) error {
	// TODO(@MattWindsor91): sanitise here
	r.Override(service.RunInfo{Args: append(armArgs[:], j.In.Litmus.Path)})
	return x.WithStdout(w).Run(ctx, r)
}

// LiftExe doesn't work.
func (Rmem) LiftExe(context.Context, backend2.LiftJob, service.RunInfo, service.Runner) error {
	return fmt.Errorf("%w: harness making", backend.ErrNotSupported)
}
