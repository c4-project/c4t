// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"

	"github.com/c4-project/c4t/internal/serviceimpl/compiler/gcc"

	"github.com/c4-project/c4t/internal/model/service"
)

// ExampleArgs is a runnable example for Args.
func ExampleArgs() {
	args := gcc.Args(
		*compiler.NewJob(compiler.Exe, nil, "a.out", "foo.c", "bar.c"),
	)
	for _, arg := range args {
		fmt.Println(arg)
	}

	// Output:
	// -o
	// a.out
	// foo.c
	// bar.c
}

// ExampleArgs_opt is a runnable example for Args that shows optimisation level selection.
func ExampleArgs_opt() {
	args := gcc.Args(
		*compiler.NewJob(
			compiler.Exe,
			&compiler.Configuration{SelectedOpt: &optlevel.Named{Name: "size"}},
			"a.out",
			"foo.c", "bar.c",
		),
	)
	for _, arg := range args {
		fmt.Println(arg)
	}

	// Output:
	// -Osize
	// -o
	// a.out
	// foo.c
	// bar.c
}

func TestArgs(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		job compiler.Job
		out []string
	}{
		"default": {
			job: *compiler.NewJob(
				compiler.Exe,
				nil,
				"a.out",
				"foo.c",
				"bar.c",
			),
			out: []string{"-o", "a.out", "foo.c", "bar.c"},
		},
		"obj": {
			job: *compiler.NewJob(
				compiler.Obj,
				nil,
				"foo.o",
				"foo.c",
			),
			out: []string{"-c", "-o", "foo.o", "foo.c"},
		},
		"with-mopt": {
			job: *compiler.NewJob(
				compiler.Exe,
				&compiler.Configuration{
					SelectedMOpt: "arch=nehalem",
				},
				"a.out",
				"foo.c",
				"bar.c",
			),
			out: []string{"-march=nehalem", "-o", "a.out", "foo.c", "bar.c"},
		},
		"with-opt": {
			job: *compiler.NewJob(
				compiler.Exe,
				&compiler.Configuration{
					SelectedOpt: &optlevel.Named{
						Name: "3",
						Level: optlevel.Level{
							Optimises:       true,
							Bias:            optlevel.BiasSpeed,
							BreaksStandards: false,
						},
					},
				},
				"a.out",
				"foo.c",
				"bar.c",
			),
			out: []string{"-O3", "-o", "a.out", "foo.c", "bar.c"},
		},
		"do-not-override-run": {
			job: *compiler.NewJob(
				compiler.Exe,
				&compiler.Configuration{
					Compiler: compiler.Compiler{
						Run: service.NewRunInfo("gcc8", "-pthread"),
					},
				},
				"a.out",
				"foo.c",
				"bar.c",
			),
			out: []string{"-o", "a.out", "foo.c", "bar.c"},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			args := gcc.Args(c.job)
			assert.Equalf(t, c.out, args, "Args(%v) didn't match", c.job)
		})
	}
}
