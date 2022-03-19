// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package plan contains the Plan type, as well as various parts of plans that don't warrant their own packages.
package plan

import (
	"context"
	"errors"
	"time"

	"github.com/c4-project/c4t/internal/timing"

	"github.com/c4-project/c4t/internal/mutation"

	backend2 "github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/subject/corpus"
)

// ErrNil is an error that can be returned if a tester stage gets a nil plan.
var ErrNil = errors.New("plan nil")

// Plan represents a test plan.
// A plan covers an entire campaign of testing.
type Plan struct {
	// Metadata contains the metadata for this plan.
	Metadata Metadata `json:"metadata"`

	// Machine represents the machine targeted by this plan.
	Machine machine.Named `json:"machine"`

	// Backend represents the backend targeted by this plan.
	Backend *backend2.NamedSpec `json:"backend,omitempty"`

	// Compilers represents the compilers to be targeted by this plan.
	Compilers compiler.InstanceMap `json:"compilers"`

	// Corpus contains each test corpus entry chosen for this plan.
	Corpus corpus.Corpus `json:"corpus"`

	// Mutation contains configuration specific to mutation testing.
	// If nil or marked disabled, no mutation testing is occurring.
	//
	// Note that the type of this field may change to an expanded form at some point.
	Mutation *mutation.Config `json:"mutation"`
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

// RunStage runs r with ctx and this plan.
// If r is a StageRunner, we marks s as completed on the resulting plan, using wall clock.
func (p *Plan) RunStage(ctx context.Context, r Runner) (*Plan, error) {
	start := time.Now()
	np, err := r.Run(ctx, p)
	if err != nil {
		return nil, err
	}
	np.Metadata.ConfirmStage(r.Stage(), timing.SpanSince(start))
	return np, nil
}

// Arches gets a list of all architectures targeted by compilers in the machine plan m.
// These architectures are in order of their string equivalents.
func (p *Plan) Arches() []id.ID {
	arches := p.archSet()
	return id.MapKeys(arches)
}

func (p *Plan) archSet() map[id.ID]struct{} {
	amap := make(map[id.ID]struct{})
	for _, c := range p.Compilers {
		amap[c.Arch] = struct{}{}
	}
	return amap
}

// CompilerIDs gets a sorted slice of all compiler IDs mentioned in this machine plan.
func (p *Plan) CompilerIDs() []id.ID {
	return id.MapKeys(p.Compilers)
}

// IsMutationTest gets whether this plan is defining a mutation test.
//
// This is shorthand for checking if the plan has an enabled mutation configuration.
func (p *Plan) IsMutationTest() bool {
	return p.Mutation != nil && p.Mutation.Enabled
}

// Mutant gets the currently selected mutant ID according to this plan, or 0 if mutation is disabled.
//
// This assumes that the plan is being repeatedly refreshed with the appropriate mutant ID.
func (p *Plan) Mutant() mutation.Mutant {
	if !p.IsMutationTest() {
		return mutation.Mutant{}
	}
	return p.Mutation.Selection
}

// SetMutant sets the current mutant ID to m, provided that mutation testing is active.
func (p *Plan) SetMutant(m mutation.Mutant) {
	if p.IsMutationTest() {
		p.Mutation.Selection = m
	}
}

// Map is shorthand for a map from machine IDs to plans.
type Map map[id.ID]Plan
