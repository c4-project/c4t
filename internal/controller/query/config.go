// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package query

import (
	"context"
	"io"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Config is the configuration for the plan query controller.
type Config struct {
	// Out is the writer to which query reports will be sent.
	Out io.Writer
}

// Run constructs a query controller from this config, then runs it.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	q, err := New(c, p)
	if err != nil {
		return nil, err
	}
	return q.Run(ctx)
}
