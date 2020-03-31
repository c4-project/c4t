// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"
)

// ExampleGCCArgs is a runnable example for GCCArgs.
func ExampleGCCArgs() {
	args := compiler.GCCArgs(service.RunInfo{
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
