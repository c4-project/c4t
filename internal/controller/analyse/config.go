// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analyse

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/saver"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Config is the configuration for the plan analyse controller.
type Config struct {
	// NWorkers is the number of parallel workers to use when performing subject analysis.
	NWorkers int
	// Observers is the list of observers to which analyses are sent.
	Observers []observer.Observer
	// SavedPaths, if present, is the pathset to which failing corpora should be sent.
	SavedPaths *saver.Pathset
}

// Run constructs a query controller from this config, then runs it.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	q, err := New(c, p)
	if err != nil {
		return nil, err
	}
	return q.Run(ctx)
}
