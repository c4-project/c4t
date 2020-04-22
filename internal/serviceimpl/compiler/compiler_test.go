// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/mock"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"
	"github.com/stretchr/testify/require"

	mdl "github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/job"
)

// TestResolver_RunCompiler tests that RunCompiler delegates properly.
func TestResolver_RunCompiler(t *testing.T) {
	var mc mock.Compiler
	r := compiler.Resolver{Compilers: map[string]compiler.Compiler{"gcc": &mc}}

	ctx := context.Background()
	j := job.Compile{
		Compiler: &mdl.Compiler{
			SelectedMOpt: "plop",
			Config: mdl.Config{
				Style: id.FromString("gcc"),
				Arch:  id.FromString("x86"),
			},
		},
		In:  []string{"foo", "bar", "baz"},
		Out: "a.out",
	}
	errw := ioutil.Discard

	mc.On("RunCompiler", ctx, j, errw).Return(nil).Once()

	err := r.RunCompiler(ctx, j, errw)
	require.NoError(t, err)
	mc.AssertExpectations(t)
}
