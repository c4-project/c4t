// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/MattWindsor91/c4t/internal/model/id"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/compiler"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/compiler/mocks"
	"github.com/stretchr/testify/require"

	mdl "github.com/MattWindsor91/c4t/internal/model/service/compiler"
)

// TestResolver_RunCompiler tests that RunCompiler delegates properly.
func TestResolver_RunCompiler(t *testing.T) {
	var mc mocks.Compiler
	r := compiler.Resolver{Compilers: map[string]compiler.Compiler{"gcc": &mc}}

	ctx := context.Background()
	j := *mdl.NewJob(
		mdl.Exe,
		&mdl.Configuration{
			SelectedMOpt: "plop",
			Compiler: mdl.Compiler{
				Style: id.FromString("gcc"),
				Arch:  id.FromString("x86"),
			},
		},
		"a.out",
		"foo", "bar", "baz",
	)
	errw := ioutil.Discard

	mc.On("RunCompiler", ctx, j, errw).Return(nil).Once()

	err := r.RunCompiler(ctx, j, errw)
	require.NoError(t, err)
	mc.AssertExpectations(t)
}
