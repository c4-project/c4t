package model

import (
	"fmt"
	"strings"
)

// ExampleParseCompilerList is a testable example for ParseCompilerlist.
func ExampleParseCompilerList() {
	list := []string{
		"localhost clang.normal gcc x86.64 enabled",
		"localhost clang.O3 gcc x86.64 enabled",
		"localhost gcc.normal gcc x86.64 enabled",
		"localhost gcc.O3 gcc x86.64 enabled",
	}
	rd := strings.NewReader(strings.Join(list, "\n"))

	compilers, err := ParseCompilerList(rd)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		for _, c := range compilers {
			fmt.Println(c.MachineId, c.Id, c.Style)
		}
	}

	// Output:
	// localhost clang.normal gcc
	// localhost clang.O3 gcc
	// localhost gcc.normal gcc
	// localhost gcc.O3 gcc
}
