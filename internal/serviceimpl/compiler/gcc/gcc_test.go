// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package gcc_test

import (
	"fmt"

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
