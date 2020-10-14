// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package resolver contains the backend resolver.
package resolver

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/model/job"
	"github.com/MattWindsor91/act-tester/internal/model/recipe"
	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/delitmus"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/herd"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/litmus"
	"github.com/MattWindsor91/act-tester/internal/subject/obs"
)

var (
	// ErrNil occurs when the backend we try to resolve is nil.
	ErrNil = errors.New("backend nil")
	// ErrUnknownStyle occurs when we ask the resolver for a backend style of which it isn't aware.
	ErrUnknownStyle = errors.New("unknown backend style")

	// Resolve is a pre-populated compiler resolver.
	Resolve = Resolver{Backends: map[string]backend.Backend{
		"delitmus": delitmus.Delitmus{},
		"herd": herdtools.Backend{
			DefaultRun: service.RunInfo{Cmd: "herd7"},
			Impl:       herd.Herd{},
		},
		"litmus": herdtools.Backend{
			DefaultRun: service.RunInfo{Cmd: "litmus7"},
			Impl:       litmus.Litmus{},
		},
	}}
)

// Resolver maps backend styles to backends.
type Resolver struct {
	// Backends is the raw map from style strings to backend runners.
	Backends map[string]backend.Backend
}

// Lift delegates lifting to the appropriate maker for j.
func (r *Resolver) Lift(ctx context.Context, j job.Lifter, errw io.Writer) (recipe.Recipe, error) {
	bi, err := r.Get(j.Backend)
	if err != nil {
		return recipe.Recipe{}, err
	}
	return bi.Lift(ctx, j, errw)
}

// ParseObs delegates observation parsing to the appropriate implementation for the backend referenced by b.
func (r *Resolver) ParseObs(ctx context.Context, b *service.Backend, rd io.Reader, o *obs.Obs) error {
	bi, err := r.Get(b)
	if err != nil {
		return err
	}
	return bi.ParseObs(ctx, b, rd, o)
}

// Get tries to look up the backend specified by b in this resolver.
func (r *Resolver) Get(b *service.Backend) (backend.Backend, error) {
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
