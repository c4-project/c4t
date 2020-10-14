// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"context"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/machine"
	"github.com/MattWindsor91/act-tester/internal/plan"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/compiler/mocks"

	"github.com/MattWindsor91/act-tester/internal/model/job/compile"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/service"
	mdl "github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler/optlevel"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/compiler"
)

// TestCompiler_Run tests running a compile job.
func TestCompiler_Run(t *testing.T) {
	var (
		mc mocks.Driver
		mp mocks.SubjectPather
	)
	mc.Test(t)
	mp.Test(t)

	names := []string{"foo", "bar", "baz"}
	c := corpus.New(names...)
	for n, cn := range c {
		err := cn.AddRecipe(id.ArchX86Skylake, recipe.New(
			n,
			recipe.AddFiles("main.c"),
			recipe.CompileAllCToExe(),
		))
		require.NoError(t, err, "setting up recipe")
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
	mc.On("RunCompiler", mock.Anything, mock.MatchedBy(func(j2 compile.Single) bool {
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
