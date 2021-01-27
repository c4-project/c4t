// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"context"
	"testing"

	mocks2 "github.com/c4-project/c4t/internal/model/service/mocks"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/serviceimpl/compiler"
	"github.com/c4-project/c4t/internal/serviceimpl/compiler/mocks"
	"github.com/stretchr/testify/require"

	mdl "github.com/c4-project/c4t/internal/model/service/compiler"
)

// TestResolver_RunCompiler tests that RunCompiler delegates properly.
func TestResolver_RunCompiler(t *testing.T) {
	mc := new(mocks.Compiler)
	mr := new(mocks2.Runner)
	mc.Test(t)
	mr.Test(t)

	r := compiler.Resolver{Compilers: map[string]compiler.Compiler{"gcc": mc}}

	ctx := context.Background()
	j := *mdl.NewJob(
		mdl.Exe,
		&mdl.Instance{
			SelectedMOpt: "plop",
			Compiler: mdl.Compiler{
				Style: id.FromString("gcc"),
				Arch:  id.FromString("x86"),
			},
		},
		"a.out",
		"foo", "bar", "baz",
	)

	mc.On("RunCompiler", ctx, j, mr).Return(nil).Once()

	err := r.RunCompiler(ctx, j, mr)
	require.NoError(t, err)
	mc.AssertExpectations(t)
	mr.AssertExpectations(t)
}
