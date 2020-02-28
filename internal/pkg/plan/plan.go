// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"context"
	"errors"
	"os"
	"sort"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

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

// Dump dumps plan p to stdout.
func (p *Plan) Dump() error {
	// TODO(@MattWindsor91): output to other files
	enc := toml.NewEncoder(os.Stdout)
	enc.Indent = "  "
	return enc.Encode(p)
}

// ParCorpus runs f for every subject in the plan's corpus.
// It threads through a context that will terminate each machine if an error occurs on some other machine.
// It also takes zero or more 'auxiliary' funcs to launch within the same context.
func (p *Plan) ParCorpus(ctx context.Context, f func(context.Context, subject.Named) error, aux ...func(context.Context) error) error {
	return p.Corpus.Par(ctx, f, aux...)
}

// Arches gets a list of all architectures targeted by compilers in the machine plan m.
// These architectures are in order of their string equivalents.
func (p *Plan) Arches() []model.ID {
	amap := make(map[string]model.ID)

	for _, c := range p.Compilers {
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
func (p *Plan) CompilerIDs() []model.ID {
	cids := make([]model.ID, len(p.Compilers))
	i := 0
	for cid := range p.Compilers {
		cids[i] = model.IDFromString(cid)
		i++
	}
	sort.Slice(cids, func(i, j int) bool {
		return cids[i].Less(cids[j])
	})
	return cids
}
