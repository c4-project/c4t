// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"sort"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// MachinePlan represents a test plan for a single machine.
type MachinePlan struct {
	// A MachinePlan subsumes a machine entry.
	model.Machine

	// Backend represents the backend targeted by this plan.
	Backend model.Backend `toml:"backend"`

	// Compilers represents the compilers to be targeted by this plan.
	// Each compiler's key is a stringified form of its machine CompilerID.
	Compilers map[string]model.Compiler `toml:"compilers"`
}

// Arches gets a list of all architectures targeted by compilers in the machine plan m.
// These architectures are in order of their string equivalents.
func (m MachinePlan) Arches() []model.ID {
	amap := make(map[string]model.ID)

	for _, c := range m.Compilers {
		amap[c.Arch.String()] = c.Arch
	}

	arches := make([]model.ID, len(amap))
	i := 0
	for _, id := range amap {
		arches[i] = id
		i++
	}

	sort.Slice(arches, func(i, j int) bool {
		return arches[i].Less(arches[j])
	})
	return arches
}

// CompilerIDs gets a sorted slice of all compiler IDs mentioned in this machine plan.
func (m MachinePlan) CompilerIDs() []model.ID {
	cids := make([]model.ID, len(m.Compilers))
	i := 0
	for cid := range m.Compilers {
		cids[i] = model.IDFromString(cid)
		i++
	}
	sort.Slice(cids, func(i, j int) bool {
		return cids[i].Less(cids[j])
	})
	return cids
}
