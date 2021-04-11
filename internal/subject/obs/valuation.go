// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import "sort"

// Valuation is an observed assignment of variable names to values.
type Valuation map[string]string

// Vars gets a sorted list of variables bound by a state.
func (v Valuation) Vars() []string {
	xs := make([]string, len(v))
	i := 0
	for x := range v {
		xs[i] = x
		i++
	}
	sort.Strings(xs)
	return xs
}
