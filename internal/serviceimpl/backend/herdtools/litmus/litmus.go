// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package litmus contains the parts of a Herdtools backend specific to litmus7.
package litmus

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/model/service/backend"

	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// Litmus describes the parts of a backend invocation that are specific to Litmus.
type Litmus struct{}

func (l Litmus) Run(ctx context.Context, j backend.LiftJob, r service.RunInfo, x service.Runner) error {
	i := Instance{Job: j, RunInfo: r, Runner: x}
	return i.Run(ctx)
}
