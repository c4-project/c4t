// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package plan contains the Plan type, as well as various parts of plans that don't warrant their own packages.
package plan

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/model/machine"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
)

// ErrNil is an error that can be returned if a tester stage gets a nil plan.
var ErrNil = errors.New("plan nil")

// Plan represents a test plan.
// A plan covers an entire campaign of testing.
type Plan struct {
	Metadata Metadata `json:"metadata"`

	// Machine represents the machine targeted by this plan.
	Machine machine.Named `json:"machine"`

	// Backend represents the backend targeted by this plan.
	Backend *service.Backend `json:"backend,omitempty"`

	// Compilers represents the compilers to be targeted by this plan.
	// Each compiler's key is a stringified form of its machine CompilerID.
	Compilers map[string]compiler.Configuration `json:"compilers"`

	// Corpus contains each test corpus entry chosen for this plan.
	Corpus corpus.Corpus `json:"corpus"`
}

// Check checks various basic properties on a plan.
func (p *Plan) Check() error {
	if err := p.Metadata.CheckVersion(); err != nil {
		return err
	}
	if len(p.Corpus) == 0 {
		return corpus.ErrNone
	}
	// TODO(@MattWindsor91): make sure compilers exist
	return nil
}

// RunStage runs f with ctx and this plan, timing it using wall-clock, and marks s as completed on the resulting plan.
func (p *Plan) RunStage(ctx context.Context, s stage.Stage, f func(context.Context, *Plan) (*Plan, error)) (*Plan, error) {
	start := time.Now()
	np, err := f(ctx, p)
	if err != nil {
		return nil, err
	}
	np.Metadata.ConfirmStage(s, start, time.Since(start))
	return np, nil
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
