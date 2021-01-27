// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import "sort"

// An observed state.
type State map[string]string

// Vars gets a sorted list of variables bound by a state.
func (s State) Vars() []string {
	vs := make([]string, len(s))
	i := 0
	for v := range s {
		vs[i] = v
		i++
	}
	sort.Strings(vs)
	return vs
}
