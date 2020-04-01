// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package gcc_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler/gcc"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// ExampleArgs is a runnable example for Args.
func ExampleArgs() {
	args := gcc.Args(service.RunInfo{
		Cmd:  "gcc7",
		Args: []string{"-funroll-loops"},
	}, job.Compile{
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
	args := gcc.Args(service.RunInfo{
		Cmd:  "gcc7",
		Args: []string{"-funroll-loops"},
	}, job.Compile{
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
