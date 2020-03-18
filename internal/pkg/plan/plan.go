// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"errors"
	"io"
	"sort"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/BurntSushi/toml"
)

// ErrNil is an error that can be returned if a tester stage gets a nil plan.
var ErrNil = errors.New("plan nil")

// plan represents a test plan.
// A plan covers an entire campaign of testing.
type Plan struct {
	Header Header `toml:"header"`

	// Machine represents the machine targeted by this plan.
	Machine model.Machine `toml:"machine"`

	// Backend represents the backend targeted by this plan.
	Backend *model.Backend `toml:"backend, omitempty"`

	// Compilers represents the compilers to be targeted by this plan.
	// Each compiler's key is a stringified form of its machine CompilerID.
	Compilers map[string]model.Compiler `toml:"compilers"`

	// Corpus contains each test corpus entry chosen for this plan.
	Corpus corpus.Corpus `toml:"corpus"`
}

// Dump dumps plan p to w.
func (p *Plan) Dump(w io.Writer) error {
	enc := toml.NewEncoder(w)
	enc.Indent = "  "
	return enc.Encode(p)
}

// Arches gets a list of all architectures targeted by compilers in the machine plan m.
// These architectures are in order of their string equivalents.
func (p *Plan) Arches() []id.ID {
	amap := make(map[string]id.ID)

	for _, c := range p.Compilers {
		amap[c.Arch.String()] = c.Arch
	}

	arches := make([]id.ID, len(amap))
	i := 0
	for _, arch := range amap {
		arches[i] = arch
		i++
	}

	sort.Slice(arches, func(i, j int) bool {
		return arches[i].Less(arches[j])
	})
	return arches
}

// CompilerIDs gets a sorted slice of all compiler IDs mentioned in this machine plan.
// It fails if any of the IDs are invalid.
func (p *Plan) CompilerIDs() ([]id.ID, error) {
	return id.MapKeys(p.Compilers)
}
