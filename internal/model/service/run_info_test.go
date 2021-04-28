// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package service_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/service"
)

// ExampleRunInfo_Override is a runnable example for the Override method.
func ExampleRunInfo_Override() {
	r1 := &service.RunInfo{
		Cmd:  "gcc",
		Args: []string{"-std=c11", "-pthread"},
		Env:  map[string]string{"FOO": "baz"},
	}
	r1.Override(service.RunInfo{
		Cmd:  "clang",
		Args: []string{"-pedantic"},
	})
	r1.Override(service.RunInfo{
		Cmd:  "",
		Args: []string{"-funroll-loops"},
		Env:  map[string]string{"FOO": "", "BAR": "baz"},
	})

	fmt.Println(r1)

	// Output:
	// BAR=baz FOO= clang -std=c11 -pthread -pedantic -funroll-loops
}

// ExampleRunInfo_Interpolate is a runnable example for the Interpolate method.
func ExampleRunInfo_Interpolate() {
	r1 := service.RunInfo{
		Cmd:  "gcc",
		Args: []string{"-std=${standard}"},
		Env:  map[string]string{"C4_MUTATION": "${mutant}"},
	}
	// Shallow copies shouldn't be affected by the interpolation
	r2 := r1
	_ = r1.Interpolate(map[string]string{"standard": "c11", "mutant": "4"})
	_ = r2.Interpolate(map[string]string{"standard": "c99", "mutant": "3"})

	fmt.Println(&r1)
	fmt.Println(&r2)

	// Output:
	// C4_MUTATION=4 gcc -std=c11
	// C4_MUTATION=3 gcc -std=c99
}

// ExampleRunInfo_SystematicID is a runnable example for the SystematicID method.
func ExampleRunInfo_SystematicID() {
	r := service.RunInfo{
		Cmd:  "gcc",
		Args: []string{"-std=c11"},
		Env:  map[string]string{"C4_MUTATION": "4"},
	}

	id1, _ := r.SystematicID()
	fmt.Println(id1)

	// If the command is empty, it doesn't appear in the ID anymore
	r.Cmd = ""
	id2, _ := r.SystematicID()
	fmt.Println(id2)

	// Output:
	// gcc.-std=c11.c4_mutation=4
	// -std=c11.c4_mutation=4
}
