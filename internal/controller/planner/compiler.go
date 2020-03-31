// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// CompilerLister is the interface of things that can query compiler information.
type CompilerLister interface {
	// ListCompilers asks the compiler inspector to list all available compilers on machine ID mid.
	ListCompilers(ctx context.Context, mid id.ID) (map[string]compiler.Config, error)
}

func (p *Planner) planCompilers(ctx context.Context) (map[string]compiler.Compiler, error) {
	cfgs, err := p.Source.CProbe.ListCompilers(ctx, p.MachineID)
	if err != nil {
		return nil, fmt.Errorf("listing compilers: %w", err)
	}

	cmps := make(map[string]compiler.Compiler, len(cfgs))
	for n, cfg := range cfgs {
		var err error
		if cmps[n], err = p.planCompiler(cfg); err != nil {
			return nil, fmt.Errorf("planning compiler %s: %w", n, err)
		}
	}

	return cmps, nil
}

func (p *Planner) planCompiler(cfg compiler.Config) (compiler.Compiler, error) {
	opt, err := p.planCompilerOpt(cfg)
	c := compiler.Compiler{
		SelectedOpt: opt,
		Config:      cfg,
	}
	return c, err
}

func (p *Planner) planCompilerOpt(_ compiler.Config) (*optlevel.Named, error) {
	// TODO(@MattWindsor91): implement this
	return nil, nil
}
