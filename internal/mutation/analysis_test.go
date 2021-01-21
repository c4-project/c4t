// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"fmt"
	"strings"

	"github.com/c4-project/c4t/internal/model/id"
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

	ana := Analysis{}
	ana.AddCompilation(compilation.Name{SubjectName: "foo", CompilerID: id.FromString("gcc")}, log, false)
	ana.AddCompilation(compilation.Name{SubjectName: "bar", CompilerID: id.FromString("clang")}, log, true)

	for mutant, hits := range ana {
		fmt.Printf("%d:", mutant)
		for _, h := range hits {
			fmt.Printf(" [%dx, %s, killed: %v]", h.NumHits, h.HitBy, h.Killed)
		}
		fmt.Println()
	}

	// Unordered output:
	// 42: [2x, foo@gcc, killed: false] [2x, bar@clang, killed: true]
	// 8: [0x, foo@gcc, killed: false] [0x, bar@clang, killed: false]
}
