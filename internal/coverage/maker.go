// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Maker contains state used by the coverage testbed maker.
type Maker struct {
}

// NewMaker constructs a new coverage testbed maker.
func NewMaker(opts ...Option) (*Maker, error) {
	m := &Maker{}
	if err := Options(opts...)(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Maker) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	// for now
	return p, nil
}
