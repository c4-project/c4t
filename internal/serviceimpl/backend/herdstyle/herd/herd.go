// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package herd contains the parts of a Herdtools backend specific to herd7.
package herd

import (
	"context"
	"fmt"
	"io"

	backend2 "github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/model/service"
)

// Herd describes the parts of a backend invocation that are specific to Herd.
type Herd struct{}

// Run fails to run Herd to generate executables.
func (h Herd) LiftExe(context.Context, backend2.LiftJob, service.RunInfo, service.Runner) error {
	return fmt.Errorf("%w: harness making", backend2.ErrNotSupported)
}

// Run runs Herd standalone.
func (h Herd) LiftStandalone(ctx context.Context, j backend2.LiftJob, r service.RunInfo, x service.Runner, w io.Writer) error {
	r.Override(service.RunInfo{Args: []string{j.In.Litmus.Path}})
	return x.WithStdout(w).Run(ctx, r)
}
