// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package plan

// MaxNumRecipes counts the upper bound on the number of recipes that need producing for this plan.
// The actual number of recipes may be lower if there is sharing between architectures (which, at time of writing,
// is not yet implemented).
func (p *Plan) MaxNumRecipes() int {
	return len(p.Arches()) * len(p.Corpus)
}

// NumExpCompilations counts the expected amount of compilations that will be produced on this plan.
// It does not actually count the number of compilations present in the plan.
func (p *Plan) NumExpCompilations() int {
	return len(p.Compilers) * len(p.Corpus)
}
