// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package compiler contains style-to-compiler resolution.
package compiler

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/c4-project/c4t/internal/stage/mach/interpreter"

	"github.com/c4-project/c4t/internal/helper/stringhelp"

	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"

	"github.com/c4-project/c4t/internal/serviceimpl/compiler/gcc"

	mdl "github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/model/service"
)

var (
	// ErrNil occurs when the compiler we try to resolve is nil.
	ErrNil = errors.New("compiler nil")
	// ErrUnknownStyle occurs when we ask the resolver for a compiler style of which it isn't aware.
	ErrUnknownStyle = errors.New("unknown compiler style")

	// CResolve is a pre-populated compiler resolver.
	CResolve = Resolver{Compilers: map[string]Compiler{
		"gcc": gcc.GCC{DefaultRun: service.RunInfo{Cmd: "gcc", Args: []string{"-pthread", "-std=gnu11"}}},
	}}
)

// Compiler contains the various interfaces that a compiler can implement.
type Compiler interface {
	mdl.Inspector
	interpreter.Driver
}

//go:generate mockery --name=Compiler

// Resolver maps compiler styles to compilers.
type Resolver struct {
	// Compilers is the raw map from style strings to compiler runners.
	Compilers map[string]Compiler
}

// Get tries to look up the compiler specified by nc in this resolver.
func (r *Resolver) Get(c *mdl.Compiler) (Compiler, error) {
	if c == nil {
		return nil, ErrNil
	}
	sstr := c.Style.String()
	cp, ok := r.Compilers[sstr]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownStyle, sstr)
	}
	return cp, nil
}

// DefaultOptLevels gets the default optimisation levels for the compiler described by c.
func (r *Resolver) DefaultOptLevels(c *mdl.Compiler) (stringhelp.Set, error) {
	cp, err := r.Get(c)
	if err != nil {
		return nil, err
	}
	return cp.DefaultOptLevels(c)
}

// OptLevels gets information about all available optimisation levels for the compiler described by c.
func (r *Resolver) OptLevels(c *mdl.Compiler) (map[string]optlevel.Level, error) {
	cp, err := r.Get(c)
	if err != nil {
		return nil, err
	}
	return cp.OptLevels(c)
}

// OptLevels gets the default machine-specific optimisation profiles for the compiler described by c.
func (r *Resolver) DefaultMOpts(c *mdl.Compiler) (stringhelp.Set, error) {
	cp, err := r.Get(c)
	if err != nil {
		return nil, err
	}
	return cp.DefaultMOpts(c)
}

// RunCompiler runs the compiler specified by nc on job j, using this resolver to map the style to a concrete compiler.
func (r *Resolver) RunCompiler(ctx context.Context, j mdl.Job, errw io.Writer) error {
	cp, err := r.Get(&j.Compiler.Compiler)
	if err != nil {
		return err
	}
	return cp.RunCompiler(ctx, j, errw)
}
