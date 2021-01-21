// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/mutation"
	"github.com/c4-project/c4t/internal/subject/compilation"
)

// ExampleStatset_Reset is a runnable example for Statset.Reset.
func ExampleStatset_Reset() {
	var s mutation.Statset
	s.Reset()

	fmt.Println("selections nil:", s.Selections == nil, "len:", len(s.Selections))
	fmt.Println("hits nil:", s.Hits == nil, "len:", len(s.Hits))
	fmt.Println("kills nil:", s.Kills == nil, "len:", len(s.Kills))

	// Output:
	// selections nil: false len: 0
	// hits nil: false len: 0
	// kills nil: false len: 0
}

// ExampleStatset_AddAnalysis is a runnable example for AddAnalysis.
func ExampleStatset_AddAnalysis() {
	var s mutation.Statset
	s.AddAnalysis(mutation.Analysis{
		27: mutation.MutantAnalysis{
			{
				NumHits: 0,
				Killed:  false,
				HitBy:   compilation.Name{SubjectName: "smooth", CompilerID: id.FromString("criminal")},
			},
			{
				NumHits: 2,
				Killed:  false,
				HitBy:   compilation.Name{SubjectName: "marco", CompilerID: id.FromString("polo")},
			},
			{
				NumHits: 3,
				Killed:  true,
				HitBy:   compilation.Name{SubjectName: "mint", CompilerID: id.FromString("polo")},
			},
		},
		53: mutation.MutantAnalysis{
			{
				NumHits: 0,
				Killed:  false,
				HitBy:   compilation.Name{SubjectName: "marco", CompilerID: id.FromString("polo")},
			},
		},
	})

	fmt.Println("27 selected:", s.Selections[27], "hit:", s.Hits[27], "killed:", s.Kills[27])
	fmt.Println("53 selected:", s.Selections[53], "hit:", s.Hits[53], "killed:", s.Kills[53])

	// Output:
	// 27 selected: 3 hit: 5 killed: 3
	// 53 selected: 1 hit: 0 killed: 0
}
