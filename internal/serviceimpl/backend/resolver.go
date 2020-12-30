// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package resolver contains the backend resolver.
package backend

import (
	"errors"
	"fmt"

	"github.com/c4-project/c4t/internal/model/id"

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

	// Resolve is a pre-populated backend resolver.
	Resolve = Resolver{Backends: map[string]backend2.Backend{
		"delitmus": delitmus.Delitmus{},
		"herdtools.herd": herdstyle.Backend{
			OptCapabilities: 0,
			Arches:          []id.ID{id.ArchC, id.ArchAArch64, id.ArchArm, id.ArchX8664, id.ArchX86, id.ArchPPC},
			DefaultRun:      service.RunInfo{Cmd: "herd7"},
			Impl:            herd.Herd{},
		},
		"herdtools.litmus": herdstyle.Backend{
			OptCapabilities: backend2.CanProduceExe,
			Arches:          []id.ID{id.ArchC, id.ArchAArch64, id.ArchArm, id.ArchX8664, id.ArchX86, id.ArchPPC},
			DefaultRun:      service.RunInfo{Cmd: "litmus7"},
			Impl:            litmus.Litmus{},
		},
		"rmem": herdstyle.Backend{
			OptCapabilities: backend2.CanLiftLitmus,
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
	Backends map[string]backend2.Backend
}

// Resolve tries to look up the backend specified by b in this resolver.
func (r *Resolver) Resolve(b *backend2.Spec) (backend2.Backend, error) {
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
