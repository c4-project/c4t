// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	backend2 "github.com/c4-project/c4t/internal/model/service/backend"
)

func (p *Planner) planBackend() (*backend2.NamedSpec, error) {
	// TODO(@MattWindsor91): add machine default backends, etc etc.
	// TODO(@MattWindsor91): don't hardcode the capability here?
	return p.source.BProbe.FindBackend(backend2.Criteria{Capability: backend2.CanLiftLitmus | backend2.CanProduceExe})
}
