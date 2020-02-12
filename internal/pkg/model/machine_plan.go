package model

import "sort"

// MachinePlan represents a test plan for a single machine.
type MachinePlan struct {
	// A MachinePlan subsumes a machine entry.
	Machine

	// Backend represents the backend targeted by this plan.
	Backend Backend `toml:"backend"`

	// Compilers represents the compilers to be targeted by this plan.
	Compilers []Compiler `toml:"compilers"`
}

// Arches gets a list of all architectures targeted by compilers in the machine plan m.
// These architectures are in order of their string equivalents.
func (m MachinePlan) Arches() []ID {
	amap := make(map[string]ID)

	for _, c := range m.Compilers {
		amap[c.Arch.String()] = c.Arch
	}

	arches := make([]ID, len(amap))
	i := 0
	for _, id := range amap {
		arches[i] = id
		i++
	}

	sort.Slice(arches, func(i, j int) bool {
		// TODO(@MattWindsor91): probably not a good heuristic
		return arches[i].String() < arches[j].String()
	})
	return arches
}
