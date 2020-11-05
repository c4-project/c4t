// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package herd contains the parts of a Herdtools backend specific to herd7.
package herd

import (
	"context"
	"fmt"

	backend2 "github.com/MattWindsor91/act-tester/internal/model/service/backend"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"

	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// Herd describes the parts of a backend invocation that are specific to Herd.
type Herd struct{}

// Run fails to run Herd (for now).
func (h Herd) Run(_ context.Context, _ backend2.LiftJob, _ service.RunInfo, _ service.Runner) error {
	// TODO(@MattWindsor91): once we extend this to deal with non-harness jobs, add functionality here.
	return fmt.Errorf("%w: harness making", backend.ErrNotSupported)
}
