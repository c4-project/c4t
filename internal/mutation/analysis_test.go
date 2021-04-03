// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation_test

import (
	"fmt"
	"strings"

	"github.com/c4-project/c4t/internal/mutation"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/subject/compilation"
)

// ExampleAnalysis_AddCompilation is a testable example for AddCompilation.
func ExampleAnalysis_AddCompilation() {
	log := strings.Join([]string{
		"warning: overfull hbox",
		"MUTATION SELECTED: 42",
		"warning: ineffective assign",
		"MUTATION HIT: 42 (barely)",
		"info: don't do this",
		"this statement is false",
		"MUTATION SELECTED: 8",
		"MUTATION HIT: 42 (somewhat)",
	}, "\n")

	ana := mutation.Analysis{}
	ana.RegisterMutant(mutation.NamedMutant(42, "XYZ", 0))

	fmt.Println("kills after 0 adds:", ana.Kills())
	ana.AddCompilation(compilation.Name{SubjectName: "foo", CompilerID: id.FromString("gcc")}, log, status.Ok)
	fmt.Println("kills after 1 adds:", ana.Kills())
	ana.AddCompilation(compilation.Name{SubjectName: "bar", CompilerID: id.FromString("clang")}, log, status.Flagged)
	fmt.Println("kills after 2 adds:", ana.Kills())

	for _, ma := range ana {
		fmt.Printf("%s:", ma.Mutant)
		for _, h := range ma.Selections {
			fmt.Printf(" [%dx, %s, killed: %v]", h.NumHits, h.HitBy, h.Killed())
		}
		fmt.Println()
	}

	// Unordered output:
	// kills after 0 adds: []
	// kills after 1 adds: []
	// kills after 2 adds: [XYZ:42]
	// XYZ:42: [2x, foo@gcc, killed: false] [2x, bar@clang, killed: true]
	// 8: [0x, foo@gcc, killed: false] [0x, bar@clang, killed: false]
}
