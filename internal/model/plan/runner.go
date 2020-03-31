// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import "context"

// Runner is the interface of parts of the tester that transform plans.
type Runner interface {
	// Run runs this type's processing stage on the plan pointed to by p.
	// It also takes a context, which can be used to cancel the process.
	// It returns an updated plan (which may or may not be p edited in-place), or an error.
	Run(ctx context.Context, p *Plan) (*Plan, error)
}
