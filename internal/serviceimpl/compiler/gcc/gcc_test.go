// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/c4t/internal/model/service/compiler"
	"github.com/MattWindsor91/c4t/internal/model/service/compiler/optlevel"

	"github.com/MattWindsor91/c4t/internal/serviceimpl/compiler/gcc"

	"github.com/MattWindsor91/c4t/internal/model/service"
)

// ExampleArgs is a runnable example for Args.
func ExampleArgs() {
	args := gcc.Args(
		*service.NewRunInfo("gcc7", "-funroll-loops"),
		*compiler.NewJob(compiler.Exe, nil, "a.out", "foo.c", "bar.c"),
	)
	for _, arg := range args {
		fmt.Println(arg)
	}

	// Output:
	// -funroll-loops
	// -o
	// a.out
	// foo.c
	// bar.c
}

// ExampleArgs_opt is a runnable example for Args that shows optimisation level selection.
func ExampleArgs_opt() {
	args := gcc.Args(
		*service.NewRunInfo("gcc7", "-funroll-loops"),
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
	// -funroll-loops
	// -Osize
	// -o
	// a.out
	// foo.c
	// bar.c
}

func TestArgs(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		run service.RunInfo
		job compiler.Job
		out []string
	}{
		"default": {
			run: *service.NewRunInfo("gcc7", "-funroll-loops"),
			job: *compiler.NewJob(
				compiler.Exe,
				nil,
				"a.out",
				"foo.c",
				"bar.c",
			),
			out: []string{"-funroll-loops", "-o", "a.out", "foo.c", "bar.c"},
		},
		"obj": {
			run: *service.NewRunInfo("gcc7", "-funroll-loops"),
			job: *compiler.NewJob(
				compiler.Obj,
				nil,
				"foo.o",
				"foo.c",
			),
			out: []string{"-funroll-loops", "-c", "-o", "foo.o", "foo.c"},
		},
		"with-mopt": {
			run: *service.NewRunInfo("gcc8"),
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
			run: *service.NewRunInfo("gcc8"),
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
			run: *service.NewRunInfo("gcc4", "-funroll-loops"),
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
			out: []string{"-funroll-loops", "-o", "a.out", "foo.c", "bar.c"},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			args := gcc.Args(c.run, c.job)
			assert.Equalf(t, c.out, args, "Args(%v, %v) didn't match", c.run, c.job)
		})
	}
}
