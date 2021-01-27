// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package interpreter_test

import (
	"context"
	"errors"
	"path"
	"testing"

	mocks2 "github.com/c4-project/c4t/internal/model/service/mocks"
	mocks3 "github.com/c4-project/c4t/internal/stage/mach/interpreter/mocks"

	"github.com/c4-project/c4t/internal/stage/mach/interpreter"

	"github.com/c4-project/c4t/internal/model/filekind"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/stretchr/testify/mock"

	mdl "github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/model/recipe"
)

// TestInterpreter_Interpret tests Interpret on an example recipe.
func TestInterpreter_Interpret(t *testing.T) {
	t.Parallel()

	mc := new(mocks3.Driver)
	mr := new(mocks2.Runner)
	mc.Test(t)
	mr.Test(t)

	r, err := recipe.New(
		"in",
		recipe.OutExe,
		recipe.AddFiles("body.c", "harness.c", "body.h"),
		recipe.CompileFileToObj(path.Join("in", "body.c")),
		recipe.CompileAllCToExe(),
	)
	require.NoError(t, err, "error while making recipe")

	c := mdl.Instance{}
	it, err := interpreter.New("a.out", r, mr, interpreter.CompileWith(mc, &c))
	require.NoError(t, err, "error while making interpreter")

	mc.On("RunCompiler",
		mock.Anything,
		*mdl.NewJob(mdl.Obj, &c, path.Join("in", "obj_0.o"), path.Join("in", "body.c")),
		mr,
	).Return(nil).Once().On("RunCompiler",
		mock.Anything,
		*mdl.NewJob(mdl.Exe, &c, "a.out", path.Join("in", "obj_0.o"), path.Join("in", "harness.c")),
		mr,
	).Return(nil).Once()

	err = it.Interpret(context.Background())
	require.NoError(t, err, "error while running interpreter")

	mc.AssertExpectations(t)
}

// TestInterpreter_Interpret_compileError tests Interpret's response to a compiler error.
func TestInterpreter_Interpret_compileError(t *testing.T) {
	t.Parallel()

	mc := new(mocks3.Driver)
	mr := new(mocks2.Runner)
	mc.Test(t)
	mr.Test(t)

	werr := errors.New("no me gusta")

	r, err := recipe.New(
		"in",
		recipe.OutExe,
		recipe.AddFiles("body.c", "harness.c", "body.h"),
		recipe.AddInstructions(recipe.Instruction{Op: recipe.Nop}),
		recipe.CompileFileToObj(path.Join("in", "body.c")),
		recipe.CompileAllCToExe(),
	)
	require.NoError(t, err, "error while making recipe")

	c := mdl.Instance{}
	it, err := interpreter.New("a.out", r, mr, interpreter.CompileWith(mc, &c))
	require.NoError(t, err, "error while making interpreter")

	mc.On("RunCompiler",
		mock.Anything,
		*mdl.NewJob(mdl.Obj, &c, path.Join("in", "obj_0.o"), path.Join("in", "body.c")),
		mr,
	).Return(werr).Once()
	// The second compile job should not be run.

	err = it.Interpret(context.Background())
	testhelp.ExpectErrorIs(t, err, werr, "wrong error while running interpreter")

	mc.AssertExpectations(t)
	mr.AssertExpectations(t)
}

// TestInterpreter_Interpret_badInstruction tests whether a bad interpreter instruction is caught correctly.
func TestInterpreter_Interpret_badInstruction(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  []recipe.Instruction
		err error
	}{
		"bad-op":   {in: []recipe.Instruction{{Op: 42}}, err: interpreter.ErrBadOp},
		"bad-file": {in: []recipe.Instruction{recipe.PushInputInst("nonsuch.c")}, err: interpreter.ErrFileUnavailable},
		"reused-file": {in: []recipe.Instruction{
			recipe.PushInputInst("body.c"),
			recipe.PushInputInst("body.c"),
		}, err: interpreter.ErrFileUnavailable,
		},
		"reused-file-inputs": {in: []recipe.Instruction{
			recipe.PushInputsInst(filekind.CSrc),
			recipe.PushInputInst("body.c"),
		}, err: interpreter.ErrFileUnavailable,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mc := new(mocks3.Driver)
			mr := new(mocks2.Runner)
			mc.Test(t)
			mr.Test(t)

			r, err := recipe.New(
				"in",
				recipe.OutExe,
				recipe.AddFiles("body.c", "harness.c", "body.h"),
				recipe.AddInstructions(c.in...),
			)
			require.NoError(t, err, "error while making recipe")

			cmp := mdl.Instance{}
			it, err := interpreter.New("a.out", r, mr, interpreter.CompileWith(mc, &cmp))
			require.NoError(t, err, "error while making interpreter")

			err = it.Interpret(context.Background())
			testhelp.ExpectErrorIs(t, err, c.err, "running interpreter on bad instruction")

			mc.AssertExpectations(t)
			mr.AssertExpectations(t)
		})
	}
}

// TestInterpreter_Interpret_tooManyObjs tests the interpreter's object overflow by setting its cap to a comically low
// amount, then overflowing it.
func TestInterpreter_Interpret_tooManyObjs(t *testing.T) {
	t.Parallel()

	mc := new(mocks3.Driver)
	mr := new(mocks2.Runner)
	mc.Test(t)
	mr.Test(t)

	r, err := recipe.New(
		"in",
		recipe.OutExe,
		recipe.AddFiles("body.c", "harness.c", "body.h"),
		recipe.CompileFileToObj(path.Join("in", "body.c")),
		recipe.CompileFileToObj(path.Join("in", "harness.c")),
	)
	require.NoError(t, err, "error while making recipe")
	c := mdl.Instance{}
	mc.On("RunCompiler",
		mock.Anything,
		*mdl.NewJob(mdl.Obj, &c, path.Join("in", "obj_0.o"), path.Join("in", "body.c")),
		mr).Return(nil).Once()

	it, err := interpreter.New("a.out", r, mr, interpreter.SetMaxObjs(1), interpreter.CompileWith(mc, &c))
	require.NoError(t, err, "error while making interpreter")

	err = it.Interpret(context.Background())
	testhelp.ExpectErrorIs(t, err, interpreter.ErrObjOverflow, "running interpreter with overflowing objs")

	mc.AssertExpectations(t)
	mr.AssertExpectations(t)
}
