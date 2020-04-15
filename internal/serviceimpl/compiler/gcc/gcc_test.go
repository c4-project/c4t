// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package gcc_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler/gcc"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// ExampleArgs is a runnable example for Args.
func ExampleArgs() {
	args := gcc.Args(
		*service.NewRunInfo("gcc7", "-funroll-loops"),
		job.Compile{
			In:  []string{"foo.c", "bar.c"},
			Out: "a.out",
		})
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
		job.Compile{
			In:       []string{"foo.c", "bar.c"},
			Out:      "a.out",
			Compiler: &compiler.Compiler{SelectedOpt: &optlevel.Named{Name: "size"}},
		})
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
		job job.Compile
		out []string
	}{
		"default": {
			run: *service.NewRunInfo("gcc7", "-funroll-loops"),
			job: job.Compile{
				In:  []string{"foo.c", "bar.c"},
				Out: "a.out",
			},
			out: []string{"-funroll-loops", "-o", "a.out", "foo.c", "bar.c"},
		},
		"with-mopt": {
			run: *service.NewRunInfo("gcc8"),
			job: job.Compile{
				Compiler: &compiler.Compiler{
					SelectedMOpt: "arch=nehalem",
				},
				In:  []string{"foo.c", "bar.c"},
				Out: "a.out",
			},
			out: []string{"-march=nehalem", "-o", "a.out", "foo.c", "bar.c"},
		},
		"with-opt": {
			run: *service.NewRunInfo("gcc8"),
			job: job.Compile{
				Compiler: &compiler.Compiler{
					SelectedOpt: &optlevel.Named{
						Name: "3",
						Level: optlevel.Level{
							Optimises:       true,
							Bias:            optlevel.BiasSpeed,
							BreaksStandards: false,
						},
					},
				},
				In:  []string{"foo.c", "bar.c"},
				Out: "a.out",
			},
			out: []string{"-O3", "-o", "a.out", "foo.c", "bar.c"},
		},
		"do-not-override-run": {
			run: *service.NewRunInfo("gcc4", "-funroll-loops"),
			job: job.Compile{
				Compiler: &compiler.Compiler{
					Config: compiler.Config{
						Run: service.NewRunInfo("gcc8", "-pthread"),
					},
				},
				In:  []string{"foo.c", "bar.c"},
				Out: "a.out",
			},
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
