// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

func (p *Planner) planBackend(ctx context.Context) (*model.Backend, error) {
	// TODO(@MattWindsor91): fix this pointer awfulness.
	return p.BProbe.FindBackend(ctx, model.IDFromString("litmus"), p.MachineID)
}
