// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package compiler contains style-to-compiler resolution.
package compiler

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler/gcc"

	mdl "github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler"
)

var (
	// ErrNil occurs when the compiler we try to resolve is nil.
	ErrNil = errors.New("compiler nil")
	// ErrUnknownStyle occurs when we ask the resolver for a compiler style of which it isn't aware.
	ErrUnknownStyle = errors.New("unknown compiler style")

	// CResolve is a pre-populated compiler resolver.
	CResolve = Resolver{Compilers: map[string]compiler.SingleRunner{
		"gcc": gcc.GCC{DefaultRun: service.RunInfo{Cmd: "gcc", Args: []string{"-pthread", "-std=gnu11"}}},
	}}
)

// Resolver maps compiler styles to compilers.
type Resolver struct {
	// Compilers is the raw map from style strings to compiler runners.
	Compilers map[string]compiler.SingleRunner
}

// Get tries to look up the compiler specified by nc in this resolver.
func (r *Resolver) Get(c *mdl.Compiler) (compiler.SingleRunner, error) {
	if c == nil {
		return nil, ErrNil
	}
	sstr := c.Style.String()
	cp, ok := r.Compilers[sstr]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownStyle, sstr)
	}
	return cp, nil
}

// RunCompiler runs the compiler specified by nc on job j, using this resolver to map the style to a concrete compiler.
func (r *Resolver) RunCompiler(ctx context.Context, j job.Compile, errw io.Writer) error {
	cp, err := r.Get(j.Compiler)
	if err != nil {
		return err
	}
	return cp.RunCompiler(ctx, j, errw)
}
