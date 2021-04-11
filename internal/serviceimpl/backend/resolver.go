// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend

import (
	"context"
	"errors"
	"fmt"

	"github.com/c4-project/c4t/internal/id"

	backend2 "github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/serviceimpl/backend/delitmus"
	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle"
	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/herd"
	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/litmus"
	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/rmem"
)

var (
	// ErrNil occurs when the backend we try to resolve is nil.
	ErrNil = errors.New("backend nil")
	// ErrUnknownStyle occurs when we ask the resolver for a backend style of which it isn't aware.
	ErrUnknownStyle = errors.New("unknown backend style")

	herdArches   = []id.ID{id.ArchC, id.ArchAArch64, id.ArchArm, id.ArchX8664, id.ArchX86, id.ArchPPC}
	litmusArches = []id.ID{id.ArchC, id.ArchAArch64, id.ArchArm, id.ArchX8664, id.ArchX86, id.ArchPPC}
	// TODO(@MattWindsor91): rmem supports more than this, but needs more work on sanitising/model selection
	rmemArches = []id.ID{id.ArchAArch64}

	// Resolve is a pre-populated backend resolver.
	Resolve = Resolver{Backends: map[string]backend2.Class{
		"delitmus": delitmus.Delitmus{},
		"herdtools.herd": herdstyle.Class{
			OptCapabilities: 0,
			Arches:          herdArches,
			DefaultRunInfo:  service.RunInfo{Cmd: "herd7"},
			Impl:            herd.Herd{},
		},
		"herdtools.litmus": herdstyle.Class{
			OptCapabilities: backend2.CanProduceExe,
			Arches:          litmusArches,
			DefaultRunInfo:  service.RunInfo{Cmd: "litmus7"},
			Impl:            litmus.Litmus{},
		},
		"rmem": herdstyle.Class{
			OptCapabilities: backend2.CanLiftLitmus,
			Arches:          rmemArches,
			DefaultRunInfo:  service.RunInfo{Cmd: "rmem"},
			Impl:            rmem.Rmem{},
		},
	}}
)

// Resolver maps backend styles to classes, and implements a resolver accordingly.
type Resolver struct {
	// Backends is the raw map from style strings to backend constructors.
	Backends map[string]backend2.Class
}

// Resolve tries to look up the backend specified by b in this resolver.
func (r *Resolver) Resolve(b backend2.Spec) (backend2.Backend, error) {
	if r == nil {
		return nil, ErrNil
	}

	sstr := b.Style.String()
	bi, ok := r.Backends[sstr]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownStyle, sstr)
	}
	return bi.Instantiate(b), nil
}

// Probe probes every class in this resolver, and aggregates the specifications.
func (r *Resolver) Probe(ctx context.Context, sr service.Runner) ([]backend2.NamedSpec, error) {
	// As an educated guess, assume every class has one spec.
	ns := make([]backend2.NamedSpec, 0, len(r.Backends))
	var (
		cns []backend2.NamedSpec
		err error
	)
	for style, c := range r.Backends {
		if cns, err = c.Probe(ctx, sr, id.FromString(style)); err != nil {
			return nil, err
		}
		ns = append(ns, cns...)
	}
	return ns, nil
}
