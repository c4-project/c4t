// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"context"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler"
	mdl "github.com/MattWindsor91/act-tester/internal/model/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/job/compile"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler/mocks"
)

// TestInterpreter_Interpret tests Interpret on an example recipe.
func TestInterpreter_Interpret(t *testing.T) {
	var mc mocks.Compiler

	r := recipe.New(
		"in",
		recipe.AddFiles("body.c", "harness.c", "body.h"),
		recipe.CompileFileToObj(path.Join("in", "body.c")),
		recipe.CompileAllCToExe(),
	)
	c := mdl.Compiler{}
	cr := compile.FromRecipe(&c, r, "a.out")
	require.ElementsMatch(t, cr.In, []string{path.Join("in", "body.c"), path.Join("in", "harness.c")},
		"filtering error making recipe")

	it, err := compiler.NewInterpreter(&mc, cr, ioutil.Discard)
	require.NoError(t, err, "error while making interpreter")

	mc.On("RunCompiler",
		mock.Anything,
		compile.New(&c, path.Join("in", "obj_0.o"), path.Join("in", "body.c")).Single(compile.Obj),
		ioutil.Discard,
	).Return(nil).Once().On("RunCompiler",
		mock.Anything,
		compile.New(&c, "a.out", path.Join("in", "obj_0.o"), path.Join("in", "harness.c")).Single(compile.Exe),
		ioutil.Discard,
	).Return(nil).Once()

	err = it.Interpret(context.Background())
	require.NoError(t, err, "error while running interpreter")

	mc.AssertExpectations(t)
}
