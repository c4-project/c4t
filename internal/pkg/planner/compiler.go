// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// CompilerLister is the interface of things that can query compiler information.
type CompilerLister interface {
	// ListCompilers asks the compiler inspector to list all available compilers given the filter f.
	ListCompilers(ctx context.Context, f model.CompilerFilter) (map[string]map[string]model.Compiler, error)
}

func (p *Planner) planCompilers(ctx context.Context) (map[string]model.Compiler, error) {
	ms := p.MachineID.String()
	flt := model.CompilerFilter{
		CompPred: p.Filter,
		// TODO(@MattWindsor91): this is singularly unpleasant!
		MachPred: fmt.Sprintf(`(id (is "%s"))`, ms),
	}

	cmap, err := p.Source.ListCompilers(ctx, flt)
	return cmap[ms], err
}
