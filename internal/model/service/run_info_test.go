// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package service_test

import (
	"fmt"

	"github.com/MattWindsor91/c4t/internal/model/service"
)

// ExampleRunInfo_Override is a runnable example for the Override method.
func ExampleRunInfo_Override() {
	r1 := service.RunInfo{
		Cmd:  "gcc",
		Args: []string{"-std=c11", "-pthread"},
	}
	r1.Override(service.RunInfo{
		Cmd:  "clang",
		Args: []string{"-pedantic"},
	})
	r1.Override(service.RunInfo{
		Cmd:  "",
		Args: []string{"-funroll-loops"},
	})

	fmt.Println("Cmd: ", r1.Cmd)
	fmt.Print("Args:")
	for _, a := range r1.Args {
		fmt.Printf(" %s", a)
	}
	fmt.Println()

	// Output:
	// Cmd:  clang
	// Args: -std=c11 -pthread -pedantic -funroll-loops
}
