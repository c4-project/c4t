// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package resolver contains the backend resolver.
package resolver

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MattWindsor91/c4t/internal/model/id"

	backend2 "github.com/MattWindsor91/c4t/internal/model/service/backend"

	"github.com/MattWindsor91/c4t/internal/model/recipe"
	"github.com/MattWindsor91/c4t/internal/model/service"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/delitmus"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle/herd"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle/litmus"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle/rmem"
	"github.com/MattWindsor91/c4t/internal/subject/obs"
)

var (
	// ErrNil occurs when the backend we try to resolve is nil.
	ErrNil = errors.New("backend nil")
	// ErrUnknownStyle occurs when we ask the resolver for a backend style of which it isn't aware.
	ErrUnknownStyle = errors.New("unknown backend style")

	// Resolve is a pre-populated backend resolver.
	Resolve = Resolver{Backends: map[string]backend.Backend{
		"delitmus": delitmus.Delitmus{},
		"herdtools.herd": herdstyle.Backend{
			OptCapabilities: 0,
			Arches:          []id.ID{id.ArchAArch64, id.ArchArm, id.ArchX8664, id.ArchX86, id.ArchPPC},
			DefaultRun:      service.RunInfo{Cmd: "herd7"},
			Impl:            herd.Herd{},
		},
		"herdtools.litmus": herdstyle.Backend{
			OptCapabilities: backend.CanProduceExe,
			Arches:          []id.ID{id.ArchC, id.ArchAArch64, id.ArchArm, id.ArchX8664, id.ArchX86, id.ArchPPC},
			DefaultRun:      service.RunInfo{Cmd: "litmus7"},
			Impl:            litmus.Litmus{},
		},
		"rmem": herdstyle.Backend{
			OptCapabilities: backend.CanLiftLitmus,
			// TODO(@MattWindsor91): rmem supports more than this, but needs more work on sanitising/model selection
			Arches:     []id.ID{id.ArchAArch64},
			DefaultRun: service.RunInfo{Cmd: "rmem"},
			Impl:       rmem.Rmem{},
		},
	}}
)

// Resolver maps backend styles to backends.
type Resolver struct {
	// Backends is the raw map from style strings to backend runners.
	Backends map[string]backend.Backend
}

// Capabilities delegates capability handling to the appropriate backend for b.
func (r *Resolver) Capabilities(b *backend2.Spec) backend.Capability {
	bi, err := r.Get(b)
	if err != nil {
		// TODO(@MattWindsor91): return something specifically stating there is no backend?
		return 0
	}
	return bi.Capabilities(b)
}

// Lift delegates lifting to the appropriate maker for j.
func (r *Resolver) Lift(ctx context.Context, j backend2.LiftJob, sr service.Runner) (recipe.Recipe, error) {
	bi, err := r.Get(j.Backend)
	if err != nil {
		return recipe.Recipe{}, err
	}
	return bi.Lift(ctx, j, sr)
}

// ParseObs delegates observation parsing to the appropriate implementation for the backend referenced by b.
func (r *Resolver) ParseObs(ctx context.Context, b *backend2.Spec, rd io.Reader, o *obs.Obs) error {
	bi, err := r.Get(b)
	if err != nil {
		return err
	}
	return bi.ParseObs(ctx, b, rd, o)
}

// Get tries to look up the backend specified by b in this resolver.
func (r *Resolver) Get(b *backend2.Spec) (backend.Backend, error) {
	if r == nil {
		return nil, ErrNil
	}

	sstr := b.Style.String()
	bi, ok := r.Backends[sstr]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownStyle, sstr)
	}
	return bi, nil
}
