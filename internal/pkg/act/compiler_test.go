// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"fmt"
	"strings"
)

// ExampleParseCompilerList is a testable example for ParseCompilerlist.
func ExampleParseCompilerList() {
	list := []string{
		"raijin clang.normal gcc x86.64 enabled",
		"fujin clang.o3 gcc x86.64 enabled",
		"raijin gcc.normal gcc x86.64 enabled",
		"fujin gcc.o3 gcc x86.64 enabled",
	}
	rd := strings.NewReader(strings.Join(list, "\n"))

	compilers, _ := ParseCompilerList(rd)

	fmt.Println(compilers["raijin"]["clang.normal"].Style)
	fmt.Println(compilers["fujin"]["gcc.o3"].Arch)

	// Output:
	// gcc
	// x86.64
}
