// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// BackendFinder is the interface of things that can find backends for machines.
type BackendFinder interface {
	// FindBackend asks for a backend with the given style on any one of machines,
	// or a default machine if none have such a backend.
	FindBackend(ctx context.Context, style model.ID, machines ...model.ID) (*model.Backend, error)
}

func (p *Planner) planBackend(ctx context.Context) (*model.Backend, error) {
	// TODO(@MattWindsor91): fix this pointer awfulness.
	return p.Source.BProbe.FindBackend(ctx, model.IDFromString("litmus"), p.MachineID)
}
