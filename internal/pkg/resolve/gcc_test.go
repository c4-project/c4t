// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package resolve_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/resolve"
)

// ExampleGCCArgs is a runnable example for GCCArgs.
func ExampleGCCArgs() {
	args := resolve.GCCArgs(model.CompilerRunInfo{
		Cmd:  "gcc7",
		Args: []string{"-funroll-loops"},
	}, model.CompileJob{
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
