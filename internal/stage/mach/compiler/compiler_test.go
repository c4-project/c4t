// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"context"
	"testing"

	mocks2 "github.com/c4-project/c4t/internal/stage/mach/interpreter/mocks"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/plan"

	"github.com/c4-project/c4t/internal/stage/mach/compiler/mocks"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"

	"github.com/c4-project/c4t/internal/subject/corpus"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/service"
	mdl "github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
	"github.com/c4-project/c4t/internal/stage/mach/compiler"
)

// TestCompiler_Run tests running a compile job.
func TestCompiler_Run(t *testing.T) {
	var (
		mc mocks2.Driver
		mp mocks.SubjectPather
	)
	mc.Test(t)
	mp.Test(t)

	names := []string{"foo", "bar", "baz"}
	c := corpus.New(names...)
	for n, cn := range c {
		r, err := recipe.New(
			n,
			recipe.OutExe,
			recipe.AddFiles("main.c"),
			recipe.CompileAllCToExe(),
		)
		require.NoError(t, err, "building recipe")
		err = cn.AddRecipe(id.ArchX86Skylake, r)
		require.NoError(t, err, "adding recipe")
		c[n] = cn
	}

	cmp := mdl.Configuration{
		SelectedMOpt: "arch=skylake",
		SelectedOpt: &optlevel.Named{
			Name: "3",
			Level: optlevel.Level{
				Optimises:       true,
				Bias:            optlevel.BiasSpeed,
				BreaksStandards: false,
			},
		},
		Compiler: mdl.Compiler{
			Style: id.CStyleGCC,
			Arch:  id.ArchX86Skylake,
			Run: &service.RunInfo{
				Cmd:  "gcc",
				Args: nil,
			},
		},
	}

	p := plan.Plan{
		Metadata: *plan.NewMetadata(0),
		Machine: machine.Named{
			ID: id.FromString("localhost"),
			Machine: machine.Machine{
				Cores: 4,
			},
		},
		Compilers: map[string]mdl.Configuration{
			"gcc": cmp,
		},
		Corpus: c,
	}

	ctx := context.Background()

	for _, n := range names {
		n := n
		mp.On("SubjectPaths", mock.MatchedBy(func(x compilation.Name) bool {
			return x.SubjectName == n
		})).Return(compilation.CompileFileset{
			Bin: "bin",
			Log: "", // disable logging
		}).Once()
	}

	// not necessarily the same context
	mc.On("RunCompiler", mock.Anything, mock.MatchedBy(func(j2 mdl.Job) bool {
		return j2.SelectedOptName() == cmp.SelectedOpt.Name && j2.SelectedMOptName() == cmp.SelectedMOpt
	}), mock.Anything).Return(nil)
	mp.On("Prepare", id.FromString("gcc")).Return(nil)

	stage, serr := compiler.New(&mc, &mp)
	require.NoError(t, serr, "constructing compile job")
	p2, err := stage.Run(ctx, &p)
	require.NoError(t, err, "running compile job")

	mp.AssertExpectations(t)
	mc.AssertExpectations(t)

	for got := range p2.Corpus {
		assert.Contains(t, names, got, "corpus got an extra subject name")
	}
}
