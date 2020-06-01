// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"context"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler/mocks"

	"github.com/MattWindsor91/act-tester/internal/model/job/compile"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/stretchr/testify/mock"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler"
	mdl "github.com/MattWindsor91/act-tester/internal/model/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// TestInstance_Compile tests running a compile job.
func TestInstance_Compile(t *testing.T) {
	var mc mocks.SingleRunner
	var mp compiler.MockSubjectPather

	names := []string{"foo", "bar", "baz"}
	c := corpus.New(names...)
	for n, cn := range c {
		err := cn.AddHarness(id.ArchX86Skylake, recipe.Recipe{
			Dir:          n,
			Files:        []string{"main.c"},
			Instructions: []recipe.Instruction{recipe.CompileBinInst()},
		})
		require.NoError(t, err, "setting up harness")
		c[n] = cn
	}

	rch := make(chan builder.Request, len(c))

	j := compiler.Instance{
		MachineID: id.FromString("localhost"),
		Compiler: &mdl.Named{
			ID: id.FromString("gcc"),
			Compiler: mdl.Compiler{
				SelectedMOpt: "arch=skylake",
				SelectedOpt: &optlevel.Named{
					Name: "3",
					Level: optlevel.Level{
						Optimises:       true,
						Bias:            optlevel.BiasSpeed,
						BreaksStandards: false,
					},
				},
				Config: mdl.Config{
					Style: id.CStyleGCC,
					Arch:  id.ArchX86Skylake,
					Run: &service.RunInfo{
						Cmd:  "gcc",
						Args: nil,
					},
				},
			},
		},
		Conf: &compiler.Config{
			Driver:    &mc,
			Observers: nil,
			Logger:    nil,
			Paths:     &mp,
			Quantities: compiler.QuantitySet{
				Timeout: 0,
			},
		},
		ResCh:  rch,
		Corpus: c,
	}

	ctx := context.Background()

	for _, n := range names {
		n := n
		mp.On("SubjectPaths", mock.MatchedBy(func(x compiler.SubjectCompile) bool {
			return x.Name == n
		})).Return(subject.CompileFileset{
			Bin: "bin",
			Log: "", // disable logging
		}).Once()
	}

	// not necessarily the same context
	mc.On("RunCompiler", mock.Anything, mock.MatchedBy(func(j2 compile.Single) bool {
		return j2.SelectedOptName() == j.Compiler.SelectedOpt.Name && j2.SelectedMOptName() == j.Compiler.SelectedMOpt
	}), mock.Anything).Return(nil)

	err := j.Compile(ctx)
	require.NoError(t, err, "running compile job")

	mp.AssertExpectations(t)
	mc.AssertExpectations(t)

	for range names {
		r := <-rch
		assert.Contains(t, names, r.Name, "builder request has weird name")
	}
}
