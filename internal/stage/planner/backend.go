// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	backend2 "github.com/MattWindsor91/c4t/internal/model/service/backend"
)

func (p *Planner) planBackend() (*backend2.NamedSpec, error) {
	// TODO(@MattWindsor91): add machine default backends, etc etc.
	return p.source.BProbe.FindBackend(backend2.Criteria{})
}
