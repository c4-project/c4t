// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stat_test

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/c4-project/c4t/internal/stat"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/mutation"
	"github.com/c4-project/c4t/internal/subject/compilation"
)

// ExampleMutation_Reset is a runnable example for Mutation.Reset.
func ExampleMutation_Reset() {
	var s stat.Mutation
	s.Reset()

	fmt.Println("by-mutant nil:", s.ByIndex == nil, "len:", len(s.ByIndex))

	// Output:
	// by-mutant nil: false len: 0
}

// ExampleMutation_AddAnalysis is a runnable example for AddAnalysis.
func ExampleMutation_AddAnalysis() {
	var s stat.Mutation
	s.AddAnalysis(mutation.Analysis{
		27: mutation.MutantAnalysis{
			Mutant: mutation.NamedMutant(27, "ABC", 1),
			Selections: []mutation.SelectionAnalysis{
				{
					NumHits: 0,
					Status:  status.Ok,
					HitBy:   compilation.Name{SubjectName: "smooth", CompilerID: id.FromString("criminal")},
				},
				{
					NumHits: 2,
					Status:  status.RunFail,
					HitBy:   compilation.Name{SubjectName: "marco", CompilerID: id.FromString("polo")},
				},
				{
					NumHits: 4,
					Status:  status.Flagged,
					HitBy:   compilation.Name{SubjectName: "mint", CompilerID: id.FromString("polo")},
				},
			},
		},
		53: mutation.MutantAnalysis{
			Mutant: mutation.NamedMutant(53, "DEF", 0),
			Selections: []mutation.SelectionAnalysis{
				{
					NumHits: 0,
					Status:  status.Filtered,
					HitBy:   compilation.Name{SubjectName: "marco", CompilerID: id.FromString("polo")},
				},
			},
		},
	})

	fmt.Println(s.ByIndex[27].Info, "selected:", s.ByIndex[27].Selections.Count, "hit:", s.ByIndex[27].Hits.Count, "killed:", s.ByIndex[27].Kills.Count)
	fmt.Println(s.ByIndex[53].Info, "selected:", s.ByIndex[53].Selections.Count, "hit:", s.ByIndex[53].Hits.Count, "killed:", s.ByIndex[53].Kills.Count)

	// Output:
	// ABC1:27 selected: 3 hit: 6 killed: 1
	// DEF:53 selected: 1 hit: 0 killed: 0
}

// ExampleMutation_DumpCSV is a runnable example for Mutation.DumpCSV.
func ExampleMutation_DumpCSV() {
	_ = (&stat.Mutation{
		ByIndex: map[mutation.Index]stat.Mutant{
			2:  {Info: mutation.AnonMutant(2), Selections: stat.Hitset{Count: 1}, Hits: stat.Hitset{Count: 0}, Kills: stat.Hitset{Count: 0}, Statuses: map[status.Status]uint64{status.Filtered: 1}},
			42: {Info: mutation.NamedMutant(42, "FOO", 0), Selections: stat.Hitset{Count: 10}, Hits: stat.Hitset{Count: 1}, Kills: stat.Hitset{Count: 0}, Statuses: map[status.Status]uint64{status.Ok: 9, status.CompileTimeout: 1}},
			53: {Info: mutation.NamedMutant(53, "BAR", 10), Selections: stat.Hitset{Count: 20}, Hits: stat.Hitset{Count: 400}, Kills: stat.Hitset{Count: 15}, Statuses: map[status.Status]uint64{status.Flagged: 15, status.CompileFail: 3, status.RunFail: 2}},
		},
	}).DumpCSV(csv.NewWriter(os.Stdout), "localhost")

	// Output:
	// localhost,2,,1,0,0,0,1,0,0,0,0,0
	// localhost,42,FOO,10,1,0,9,0,0,0,1,0,0
	// localhost,53,BAR10,20,400,15,0,0,15,3,0,2,0
}
